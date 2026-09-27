package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS configures cross-origin resource sharing.
// Supports configured origins and automatically permits localhost/127.0.0.1 in development.
func CORS(allowedOrigins []string) fiber.Handler {
	allowedMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		allowedMap[strings.TrimSpace(o)] = true
	}

	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			if allowedMap[origin] {
				return true
			}
			// Automatically permit localhost and 127.0.0.1 on development ports
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				origin == "http://localhost" ||
				origin == "http://127.0.0.1" {
				return true
			}
			return false
		},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Branch-Id",
		},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
	})
}
