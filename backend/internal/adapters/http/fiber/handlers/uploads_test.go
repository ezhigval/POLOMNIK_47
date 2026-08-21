package handlers_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/adapters/http/fiber/handlers"
	"polomnik/internal/config"
)

func TestUploadImageStoresPNGAndReturnsURL(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := config.Config{
		UploadDir:           dir,
		UploadPublicBaseURL: "http://localhost:8080",
		UploadMaxBytes:      2 * 1024 * 1024,
	}

	app := fiber.New(fiber.Config{BodyLimit: 3 * 1024 * 1024})
	app.Post("/uploads", handlers.UploadImage(cfg))
	app.Static("/uploads", dir)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "cover.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	if err := png.Encode(part, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req, err := http.NewRequest(fiber.MethodPost, "/uploads", &body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		payload, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, payload)
	}

	var envelope dto.DataEnvelope[dto.UploadResponse]
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if envelope.Data.URL == "" || envelope.Data.Path == "" {
		t.Fatalf("empty upload response: %+v", envelope.Data)
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(envelope.Data.Path))); err != nil {
		t.Fatalf("stored file missing: %v", err)
	}

	getReq, err := http.NewRequest(fiber.MethodGet, "/uploads/"+envelope.Data.Path, nil)
	if err != nil {
		t.Fatalf("new get request: %v", err)
	}
	getResp, err := app.Test(getReq, -1)
	if err != nil {
		t.Fatalf("static get: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected static 200, got %d", getResp.StatusCode)
	}
}

func TestUploadImageRejectsNonImage(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := config.Config{
		UploadDir:           dir,
		UploadPublicBaseURL: "http://localhost:8080",
		UploadMaxBytes:      1024 * 1024,
	}

	app := fiber.New()
	app.Post("/uploads", handlers.UploadImage(cfg))

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "notes.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := io.WriteString(part, "not an image"); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req, err := http.NewRequest(fiber.MethodPost, "/uploads", &body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}
