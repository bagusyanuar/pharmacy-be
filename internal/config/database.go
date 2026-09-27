package config

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ConnectDatabase opens the database connection pool based on configuration.
func ConnectDatabase(cfg *Config, log *zap.Logger) *gorm.DB {
	dsn := cfg.DSN()

	gormLogLevel := gormlogger.Silent
	if cfg.AppEnv == "development" {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get generic database object", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}

	log.Info("database connected successfully", zap.String("database", cfg.DBName))
	return db
}
