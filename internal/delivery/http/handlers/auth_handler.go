package handlers

import (
	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(usecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// Login godoc
// @Summary Admin login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}
	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}
	return response.JSON(c, fiber.StatusOK, fiber.Map{"access_token": token})
}
