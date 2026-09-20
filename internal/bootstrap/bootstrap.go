package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/bagusyanuar/pharmacy-be/internal/config"
	"github.com/bagusyanuar/pharmacy-be/internal/middleware"
	authdelivery "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/delivery"
	authusecase "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/usecase"
	userdelivery "github.com/bagusyanuar/pharmacy-be/internal/modules/user/delivery"
	userdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	userinfra "github.com/bagusyanuar/pharmacy-be/internal/modules/user/infrastructure"
	userusecase "github.com/bagusyanuar/pharmacy-be/internal/modules/user/usecase"
	"github.com/bagusyanuar/pharmacy-be/pkg/jwt"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// SetupApp initializes and configures the Fiber application with all modules and middleware.
func SetupApp(cfg *config.Config, db *gorm.DB, log *zap.Logger) *fiber.App {
	// Auto migrate tables (useful in development)
	if cfg.AppEnv == "development" {
		if err := db.AutoMigrate(&userdomain.User{}); err != nil {
			log.Warn("auto-migration warning", zap.Error(err))
		}
	}

	app := fiber.New(fiber.Config{
		AppName: "Go Backend Template API",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(middleware.ZapLogger(log))
	app.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "service is healthy", fiber.Map{
			"status": "ok",
			"env":    cfg.AppEnv,
		})
	})

	// Base API route group
	api := app.Group("/api/v1")

	// Core utilities
	accessManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiration, cfg.JWTIssuer)
	refreshManager := jwt.NewManager(cfg.JWTRefreshSecret, cfg.JWTRefreshExpiration, cfg.JWTIssuer)
	authMiddleware := middleware.Auth(accessManager)

	// Module: User
	userRepo := userinfra.NewUserRepository(db)
	userUC := userusecase.NewUserUsecase(userRepo, log)
	userdelivery.NewUserHandler(api, userUC, authMiddleware)

	// Module: Auth
	authUC := authusecase.NewAuthUsecase(userRepo, accessManager, refreshManager, log)
	authdelivery.NewAuthHandler(
		api,
		authUC,
		authMiddleware,
		cfg.RefreshTokenCookieName,
		cfg.AppEnv != "development",
		cfg.JWTRefreshExpiration,
	)

	return app
}
