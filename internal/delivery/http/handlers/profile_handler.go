package handlers

import (
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	usecase *usecase.ProfileUsecase
}

type upsertProfileRequest struct {
	FullName    string `json:"full_name"`
	Headline    string `json:"headline"`
	Bio         string `json:"bio"`
	GitHubURL   string `json:"github_url"`
	LinkedInURL string `json:"linkedin_url"`
	AvatarURL   string `json:"avatar_url"`
}

func NewProfileHandler(usecase *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{usecase: usecase}
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	profile, err := h.usecase.Get()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if profile == nil {
		return response.Error(c, fiber.StatusNotFound, "profile not found")
	}
	return response.JSON(c, fiber.StatusOK, profile)
}

func (h *ProfileHandler) Upsert(c *fiber.Ctx) error {
	var req upsertProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	profile := &entities.Profile{
		FullName:    req.FullName,
		Headline:    req.Headline,
		Bio:         req.Bio,
		GitHubURL:   req.GitHubURL,
		LinkedInURL: req.LinkedInURL,
		AvatarURL:   req.AvatarURL,
	}
	if err := h.usecase.Upsert(profile); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.JSON(c, fiber.StatusOK, fiber.Map{"message": "profile saved"})
}
