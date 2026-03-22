package handlers

import (
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	usecase *usecase.ProjectUsecase
}

func NewProjectHandler(usecase *usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{usecase: usecase}
}

func (h *ProjectHandler) ListFeatured(c *fiber.Ctx) error {
	projects, err := h.usecase.ListFeatured()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.JSON(c, fiber.StatusOK, projects)
}
