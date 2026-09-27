package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/bagusyanuar/pharmacy-be/internal/bootstrap"
	"github.com/bagusyanuar/pharmacy-be/internal/config"
	applogger "github.com/bagusyanuar/pharmacy-be/pkg/logger"
)

// @title           Pharmacy Information System (POS & ERP) API
// @version         1.0
// @description     High-performance multi-branch pharmacy management REST API built with Go, Fiber v3, GORM, and PostgreSQL.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Pharmacy Engineering
// @contact.url     https://github.com/bagusyanuar/pharmacy-be

// @license.name    Proprietary
// @license.url     https://github.com/bagusyanuar/pharmacy-be

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and your JWT token.

func main() {
	cfg := config.Load()

	log := applogger.Init(applogger.Config{
		Env:        cfg.AppEnv,
		FilePath:   cfg.LogFilePath,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	})
	defer applogger.Sync()

	db := config.ConnectDatabase(cfg, log)

	app := bootstrap.SetupApp(cfg, db, log)

	// Graceful shutdown hooks
	app.Hooks().OnPreShutdown(func() error {
		log.Info("shutdown signal received, draining in-flight requests")
		return nil
	})
	app.Hooks().OnPostShutdown(func(err error) error {
		if err != nil {
			log.Error("graceful shutdown finished with error", zap.Error(err))
		}
		if sqlDB, dbErr := db.DB(); dbErr == nil {
			if closeErr := sqlDB.Close(); closeErr != nil {
				log.Error("failed to close database connection", zap.Error(closeErr))
			} else {
				log.Info("database connection closed")
			}
		}
		return nil
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info("starting HTTP server", zap.String("port", cfg.AppPort), zap.String("env", cfg.AppEnv))

	err := app.Listen(":"+cfg.AppPort, fiber.ListenConfig{
		GracefulContext: ctx,
		ShutdownTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Error("server stopped unexpectedly", zap.Error(err))
	} else {
		log.Info("server shut down gracefully")
	}
}
