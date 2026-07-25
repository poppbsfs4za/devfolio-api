package middleware

import (
	"github.com/example/devfolio-api/internal/database"
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

// RequireDB short-circuits DB-dependent routes with a structured 503 when no
// database connection is currently available (e.g. it failed at startup
// because Neon's compute/connection quota was exhausted). This keeps the
// process alive and serving /health and /ready(z) instead of the whole app
// crashing or returning confusing 500s from nil-pointer access deeper in the
// repository layer.
func RequireDB(status *database.Status) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !status.Available() {
			return response.JSON(c, fiber.StatusServiceUnavailable, fiber.Map{
				"status": "unavailable",
				"error": fiber.Map{
					"code":    "DB_UNAVAILABLE",
					"message": "service temporarily unavailable: database is not connected",
				},
			})
		}
		return c.Next()
	}
}
