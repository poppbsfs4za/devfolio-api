package response

import "github.com/gofiber/fiber/v2"

type successResponse struct {
	Data any `json:"data"`
}

type errorBody struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error any `json:"error"`
}

func JSON(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(successResponse{
		Data: data,
	})
}

// Backward-compatible:
// - response.Error(c, 400, "invalid request")
// - response.Error(c, 400, "BAD_REQUEST", "invalid request")
func Error(c *fiber.Ctx, status int, args ...string) error {
	switch len(args) {
	case 0:
		return c.Status(status).JSON(errorResponse{
			Error: errorBody{
				Code:    "UNKNOWN_ERROR",
				Message: "unknown error",
			},
		})
	case 1:
		return c.Status(status).JSON(errorResponse{
			Error: errorBody{
				Message: args[0],
			},
		})
	default:
		return c.Status(status).JSON(errorResponse{
			Error: errorBody{
				Code:    args[0],
				Message: args[1],
			},
		})
	}
}
