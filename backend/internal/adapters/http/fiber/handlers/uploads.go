package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/config"
)

var allowedImageMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func UploadImage(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusUnprocessableEntity,
				Code:    "VALIDATION_ERROR",
				Message: "Expected multipart field \"file\"",
			})
		}

		if fileHeader.Size <= 0 {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusUnprocessableEntity,
				Code:    "VALIDATION_ERROR",
				Message: "Empty file",
			})
		}
		if cfg.UploadMaxBytes > 0 && fileHeader.Size > int64(cfg.UploadMaxBytes) {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusRequestEntityTooLarge,
				Code:    "FILE_TOO_LARGE",
				Message: fmt.Sprintf("File exceeds max size of %d bytes", cfg.UploadMaxBytes),
			})
		}

		src, err := fileHeader.Open()
		if err != nil {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to open uploaded file",
			})
		}
		defer src.Close()

		head := make([]byte, 512)
		n, err := io.ReadFull(src, head)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to read uploaded file",
			})
		}
		head = head[:n]

		contentType := http.DetectContentType(head)
		if declared := strings.TrimSpace(fileHeader.Header.Get("Content-Type")); declared != "" {
			if mediaType, _, err := mime.ParseMediaType(declared); err == nil {
				if _, ok := allowedImageMIME[mediaType]; ok && mediaType != "application/octet-stream" {
					// Prefer sniff when possible; fall back to declared only if sniff is generic.
					if contentType == "application/octet-stream" || contentType == "text/plain; charset=utf-8" {
						contentType = mediaType
					}
				}
			}
		}

		ext, ok := allowedImageMIME[contentType]
		if !ok {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusUnprocessableEntity,
				Code:    "UNSUPPORTED_MEDIA_TYPE",
				Message: "Only JPEG, PNG, WebP and GIF images are allowed",
			})
		}

		now := time.Now().UTC()
		relDir := filepath.Join(fmt.Sprintf("%04d", now.Year()), fmt.Sprintf("%02d", int(now.Month())))
		relPath := filepath.ToSlash(filepath.Join(relDir, uuid.NewString()+ext))
		absDir := filepath.Join(cfg.UploadDir, relDir)
		if err := os.MkdirAll(absDir, 0o755); err != nil {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create upload directory",
			})
		}

		absPath := filepath.Join(cfg.UploadDir, filepath.FromSlash(relPath))
		dst, err := os.OpenFile(absPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to store uploaded file",
			})
		}
		defer dst.Close()

		if _, err := dst.Write(head); err != nil {
			_ = os.Remove(absPath)
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to store uploaded file",
			})
		}
		if _, err := io.Copy(dst, src); err != nil {
			_ = os.Remove(absPath)
			return writeAppError(c, &AppError{
				Status:  fiber.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Failed to store uploaded file",
			})
		}

		base := strings.TrimRight(cfg.UploadPublicBaseURL, "/")
		url := base + "/uploads/" + relPath

		return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.UploadResponse]{
			Data: dto.UploadResponse{
				URL:  url,
				Path: relPath,
			},
		})
	}
}
