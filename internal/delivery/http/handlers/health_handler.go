package handlers

import (
	"log"

	"github.com/example/devfolio-api/internal/database"
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	dbStatus *database.Status
}

func NewHealthHandler(dbStatus *database.Status) *HealthHandler {
	return &HealthHandler{dbStatus: dbStatus}
}

// Health godoc
// @Summary Liveness check
// @Description Reports whether the API process is up. Does not touch the database.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ok"})
}

// Ready godoc
// @Summary Readiness check
// @Description Reports whether the API can currently reach the database. Returns 503 if not.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	if err := h.dbStatus.Ping(c.Context()); err != nil {
		log.Printf("[readyz] database not reachable: %v", err)
		return response.JSON(c, fiber.StatusServiceUnavailable, fiber.Map{
			"status": "unavailable",
			"error": fiber.Map{
				"code":    "DB_UNAVAILABLE",
				"message": "database is not reachable",
			},
		})
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ready"})
}
