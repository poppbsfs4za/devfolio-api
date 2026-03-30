package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MaxCoverUploadSize int64 = 5 * 1024 * 1024 // 5MB

type UploadUsecase struct {
	uploadDir     string
	publicBaseURL string
}

func NewUploadUsecase(uploadDir string, publicBaseURL string) *UploadUsecase {
	return &UploadUsecase{
		uploadDir:     uploadDir,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}
}

type SaveCoverInput struct {
	File io.ReadSeeker
	Size int64
}

type SaveCoverOutput struct {
	URL      string
	Filename string
	Path     string
}

func (u *UploadUsecase) SaveCover(input SaveCoverInput) (*SaveCoverOutput, error) {
	if input.File == nil {
		return nil, errors.New("file is required")
	}

	if input.Size > MaxCoverUploadSize {
		return nil, errors.New("file size must be less than or equal to 5 MB")
	}

	// sniff file type
	sniff := make([]byte, 512)
	n, err := input.File.Read(sniff)
	if err != nil && err != io.EOF {
		return nil, errors.New("failed to read uploaded file")
	}

	contentType := http.DetectContentType(sniff[:n])

	ext, ok := allowedImageExtension(contentType)
	if !ok {
		return nil, errors.New("only jpg, jpeg, png, webp, and gif files are allowed")
	}

	// reset cursor
	if _, err := input.File.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("failed to process uploaded file")
	}

	// ensure dir
	if err := os.MkdirAll(u.uploadDir, 0755); err != nil {
		return nil, errors.New("failed to prepare upload directory")
	}

	filename := buildUploadFilename(ext)
	dstPath := filepath.Join(u.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, errors.New("failed to create upload file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, input.File); err != nil {
		return nil, errors.New("failed to save uploaded file")
	}

	fileURL := fmt.Sprintf("%s/uploads/covers/%s", u.publicBaseURL, filename)

	return &SaveCoverOutput{
		URL:      fileURL,
		Filename: filename,
		Path:     dstPath,
	}, nil
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
