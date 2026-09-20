package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	AppEnv  string

	// Database
	DBDriver   string
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBTimeZone string

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	// Logging
	LogFilePath   string
	LogMaxSizeMB  int
	LogMaxBackups int
	LogMaxAgeDays int
	LogCompress   bool

	// Auth & JWT
	JWTSecret     string
	JWTExpiration time.Duration
	JWTIssuer     string

	JWTRefreshSecret     string
	JWTRefreshExpiration time.Duration

	RefreshTokenCookieName string

	// CORS
	CORSAllowedOrigins []string
}

func Load() *Config {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("APP_PORT", "3000")
	v.SetDefault("APP_ENV", "development")

	v.SetDefault("DB_DRIVER", "postgres")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "app_db")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "UTC")

	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "30m")
	v.SetDefault("DB_CONN_MAX_IDLE_TIME", "5m")

	v.SetDefault("LOG_FILE_PATH", "storage/logs/app.log")
	v.SetDefault("LOG_MAX_SIZE_MB", 100)
	v.SetDefault("LOG_MAX_BACKUPS", 7)
	v.SetDefault("LOG_MAX_AGE_DAYS", 30)
	v.SetDefault("LOG_COMPRESS", true)

	v.SetDefault("JWT_SECRET", "super-secret-jwt-key-change-in-production")
	v.SetDefault("JWT_EXPIRATION", "24h")
	v.SetDefault("JWT_ISSUER", "go-backend-template")

	v.SetDefault("JWT_REFRESH_SECRET", "super-secret-refresh-key-change-in-production")
	v.SetDefault("JWT_REFRESH_EXPIRATION", "168h")

	v.SetDefault("REFRESH_TOKEN_COOKIE_NAME", "refresh_token")

	if err := v.ReadInConfig(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	return &Config{
		AppPort: v.GetString("APP_PORT"),
		AppEnv:  v.GetString("APP_ENV"),

		DBDriver:   v.GetString("DB_DRIVER"),
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetInt("DB_PORT"),
		DBName:     v.GetString("DB_NAME"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),
		DBTimeZone: v.GetString("DB_TIMEZONE"),

		DBMaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		DBMaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		DBConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		DBConnMaxIdleTime: v.GetDuration("DB_CONN_MAX_IDLE_TIME"),

		LogFilePath:   v.GetString("LOG_FILE_PATH"),
		LogMaxSizeMB:  v.GetInt("LOG_MAX_SIZE_MB"),
		LogMaxBackups: v.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAgeDays: v.GetInt("LOG_MAX_AGE_DAYS"),
		LogCompress:   v.GetBool("LOG_COMPRESS"),

		JWTSecret:     v.GetString("JWT_SECRET"),
		JWTExpiration: v.GetDuration("JWT_EXPIRATION"),
		JWTIssuer:     v.GetString("JWT_ISSUER"),

		JWTRefreshSecret:     v.GetString("JWT_REFRESH_SECRET"),
		JWTRefreshExpiration: v.GetDuration("JWT_REFRESH_EXPIRATION"),

		RefreshTokenCookieName: v.GetString("REFRESH_TOKEN_COOKIE_NAME"),

		CORSAllowedOrigins: parseCommaList(v.GetString("CORS_ALLOWED_ORIGINS")),
	}
}

func parseCommaList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
