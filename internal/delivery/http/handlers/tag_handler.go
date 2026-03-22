package handlers

import (
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	usecase *usecase.TagUsecase
}

func NewTagHandler(usecase *usecase.TagUsecase) *TagHandler {
	return &TagHandler{usecase: usecase}
}

func (h *TagHandler) List(c *fiber.Ctx) error {
	tags, err := h.usecase.List()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.JSON(c, fiber.StatusOK, tags)
}

func (h *TagHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	tag, err := h.usecase.Create(req.Name)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.JSON(c, fiber.StatusCreated, tag)
}
