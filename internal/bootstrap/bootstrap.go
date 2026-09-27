package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/bagusyanuar/pharmacy-be/docs"
	"github.com/bagusyanuar/pharmacy-be/internal/config"
	"github.com/bagusyanuar/pharmacy-be/internal/middleware"
	authdelivery "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/delivery"
	authusecase "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/usecase"
	userdelivery "github.com/bagusyanuar/pharmacy-be/internal/modules/user/delivery"
	userinfra "github.com/bagusyanuar/pharmacy-be/internal/modules/user/infrastructure"
	userusecase "github.com/bagusyanuar/pharmacy-be/internal/modules/user/usecase"
	"github.com/bagusyanuar/pharmacy-be/pkg/jwt"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// SetupApp initializes and configures the Fiber application with all modules and middleware.
func SetupApp(cfg *config.Config, db *gorm.DB, log *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Pharmacy Information System API",
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

	// Swagger documentation specifications & interactive UI (embedded in binary)
	app.Get("/docs/swagger.yaml", func(c fiber.Ctx) error {
		content, err := docs.FS.ReadFile("swagger.yaml")
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("documentation specification not found")
		}
		c.Set(fiber.HeaderContentType, "application/yaml")
		return c.Send(content)
	})
	app.Get("/docs/modules/:name", func(c fiber.Ctx) error {
		content, err := docs.FS.ReadFile("modules/" + c.Params("name"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("module specification not found")
		}
		c.Set(fiber.HeaderContentType, "application/yaml")
		return c.Send(content)
	})
	app.Get("/docs", func(c fiber.Ctx) error {
		content, err := docs.FS.ReadFile("index.html")
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("documentation page not found")
		}
		c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
		return c.Send(content)
	})
	app.Get("/swagger", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusFound).To("/docs")
	})
	app.Get("/swagger/*", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusFound).To("/docs")
	})
	app.Get("/api-docs", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusFound).To("/docs")
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
