package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

const MaxCoverUploadSize int64 = 5 * 1024 * 1024 // 5MB

type UploadUsecase struct {
	bucketName    string
	publicBaseURL string
	storageClient *storage.Client
}

func NewUploadUsecase(bucketName string, publicBaseURL string, storageClient *storage.Client) *UploadUsecase {
	return &UploadUsecase{
		bucketName:    bucketName,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
		storageClient: storageClient,
	}
}

type SaveCoverInput struct {
	File io.ReadSeeker
	Size int64
}

type SaveCoverOutput struct {
	URL       string
	Filename  string
	ObjectKey string
}

func (u *UploadUsecase) SaveCover(ctx context.Context, input SaveCoverInput) (*SaveCoverOutput, error) {
	if input.File == nil {
		return nil, errors.New("file is required")
	}

	if input.Size > MaxCoverUploadSize {
		return nil, errors.New("file size must be less than or equal to 5 MB")
	}

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

	if _, err := input.File.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("failed to process uploaded file")
	}

	filename := buildUploadFilename(ext)
	objectKey := fmt.Sprintf("covers/%s", filename)

	writer := u.storageClient.
		Bucket(u.bucketName).
		Object(objectKey).
		NewWriter(ctx)

	writer.ContentType = contentType
	writer.CacheControl = "public, max-age=31536000"

	if _, err := io.Copy(writer, input.File); err != nil {
		_ = writer.Close()
		return nil, errors.New("failed to save uploaded file")
	}

	if err := writer.Close(); err != nil {
		return nil, errors.New("failed to finalize uploaded file")
	}

	fileURL := fmt.Sprintf("%s/%s", u.publicBaseURL, objectKey)

	return &SaveCoverOutput{
		URL:       fileURL,
		Filename:  filename,
		ObjectKey: objectKey,
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
