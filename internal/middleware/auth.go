package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/bagusyanuar/pharmacy-be/pkg/jwt"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// Auth guards route groups requiring a valid Bearer token.
func Auth(manager *jwt.Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return response.Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
		}

		claims, err := manager.ParseToken(parts[1])
		if err != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
		}

		c.Locals(string(UserIDContextKey), claims.UserID)
		return c.Next()
	}
}

// GetAuthUserID retrieves the authenticated user ID from context.
func GetAuthUserID(c fiber.Ctx) string {
	val, ok := c.Locals(string(UserIDContextKey)).(string)
	if !ok {
		return ""
	}
	return val
}
