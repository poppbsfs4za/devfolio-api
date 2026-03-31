package handlers

import (
	"time"

	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

const authCookieName = "devfolio_token"

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

type loginRequest struct {
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
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid email or password")
	}

	c.Cookie(buildAuthCookie(token, time.Now().Add(24*time.Hour), true))

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"access_token": token,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {

	c.Cookie(buildAuthCookie("", time.Now().Add(-1*time.Hour), true))

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "logged out",
	})
}

func buildAuthCookie(value string, expires time.Time, secure bool) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     authCookieName,
		Value:    value,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   secure,
		Expires:  expires,
	}
}
