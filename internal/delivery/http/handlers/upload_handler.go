package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

const maxCoverUploadSize = 5 * 1024 * 1024 // 5 MB

type UploadHandler struct {
	uploadDir     string
	publicBaseURL string
}

func NewUploadHandler(uploadDir string, publicBaseURL string) *UploadHandler {
	return &UploadHandler{
		uploadDir:     uploadDir,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}
}

func (h *UploadHandler) UploadCover(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file is required")
	}

	if fileHeader.Size > maxCoverUploadSize {
		return response.Error(c, fiber.StatusBadRequest, "file size must be less than or equal to 5 MB")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "failed to open uploaded file")
	}
	defer src.Close()

	sniff := make([]byte, 512)
	n, err := src.Read(sniff)
	if err != nil && err != io.EOF {
		return response.Error(c, fiber.StatusBadRequest, "failed to read uploaded file")
	}

	contentType := http.DetectContentType(sniff[:n])
	ext, ok := allowedImageExtension(contentType)
	if !ok {
		return response.Error(c, fiber.StatusBadRequest, "only jpg, jpeg, png, webp, and gif files are allowed")
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to process uploaded file")
	}

	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to prepare upload directory")
	}

	filename := buildUploadFilename(ext)
	dstPath := filepath.Join(h.uploadDir, filename)

	fmt.Println("saving upload to:", dstPath)

	dst, err := os.Create(dstPath)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create upload file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to save uploaded file")
	}

	fileURL := fmt.Sprintf("%s/uploads/covers/%s", h.publicBaseURL, filename)

	return response.JSON(c, fiber.StatusCreated, fiber.Map{
		"url":      fileURL,
		"filename": filename,
	})
}

func allowedImageExtension(contentType string) (string, bool) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/webp":
		return ".webp", true
	case "image/gif":
		return ".gif", true
	default:
		return "", false
	}
}

func buildUploadFilename(ext string) string {
	randomBytes := make([]byte, 8)
	_, _ = rand.Read(randomBytes)

	return fmt.Sprintf(
		"cover_%d_%s%s",
		time.Now().UnixNano(),
		hex.EncodeToString(randomBytes),
		ext,
	)
}
