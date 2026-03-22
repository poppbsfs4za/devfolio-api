package middleware

import (
	"strings"

	"github.com/example/devfolio-api/internal/delivery/http/response"
	pkgAuth "github.com/example/devfolio-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

func JWT(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, fiber.StatusUnauthorized, "invalid authorization header")
		}
		claims, err := pkgAuth.ParseToken(secret, parts[1])
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("display_name", claims.DisplayName)
		return c.Next()
	}
}
