package handlers

import (
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// Health godoc
// @Summary Health check
// @Description Check API status
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ok"})
}
