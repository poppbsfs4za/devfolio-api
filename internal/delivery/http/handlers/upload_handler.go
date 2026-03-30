package handlers

import (
	"os"

	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	usecase *usecase.UploadUsecase
}

func NewUploadHandler(usecase *usecase.UploadUsecase) *UploadHandler {
	return &UploadHandler{usecase: usecase}
}

func (h *UploadHandler) UploadCover(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "file is required")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "UPLOAD_OPEN_FAILED", "failed to open uploaded file")
	}
	defer src.Close()

	readSeeker, ok := src.(interface {
		Read([]byte) (int, error)
		Seek(int64, int) (int64, error)
	})
	if !ok {
		return response.Error(c, fiber.StatusInternalServerError, "UPLOAD_INVALID_STREAM", "uploaded file stream is not seekable")
	}

	result, err := h.usecase.SaveCover(usecase.SaveCoverInput{
		File: readSeeker,
		Size: fileHeader.Size,
	})
	if err != nil {
		switch err.Error() {
		case "file size must be less than or equal to 5 MB":
			return response.Error(c, fiber.StatusBadRequest, "UPLOAD_TOO_LARGE", err.Error())
		case "only jpg, jpeg, png, webp, and gif files are allowed":
			return response.Error(c, fiber.StatusBadRequest, "UPLOAD_INVALID_FILE", err.Error())
		default:
			return response.Error(c, fiber.StatusBadRequest, "UPLOAD_FAILED", err.Error())
		}
	}

	if _, err := os.Stat(result.Path); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "UPLOAD_SAVE_VERIFICATION_FAILED", "uploaded file was not saved correctly")
	}

	return response.JSON(c, fiber.StatusCreated, fiber.Map{
		"url":      result.URL,
		"filename": result.Filename,
	})
}
