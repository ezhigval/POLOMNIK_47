package fiberhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"palomnik/internal/adapters/repository/memory"
)

func TestRegisterRequiresPersonalDataConsent(t *testing.T) {
	app := newTestAppWithStore(memoryStore(t), "secret-token")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{
		"name": "Тест",
		"email": "consent-reg@example.com",
		"phone": "+79001112233",
		"password": "password1"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 without consent, got %d", resp.StatusCode)
	}
}

func TestRegisterRecordsConsent(t *testing.T) {
	store := memoryStore(t)
	app := newTestAppWithStore(store, "secret-token")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{
		"name": "Тест",
		"email": "consent-ok@example.com",
		"phone": "+79001112234",
		"password": "password1",
		"consent_personal_data": true,
		"consent_marketing": true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Token string `json:"token"`
			User  struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/me/consents", nil)
	listReq.Header.Set("Authorization", "Bearer "+body.Data.Token)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list consents: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}
	var listBody struct {
		Data []struct {
			ConsentType     string `json:"consent_type"`
			DocumentVersion string `json:"document_version"`
		} `json:"data"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&listBody); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listBody.Data) < 2 {
		t.Fatalf("expected personal_data and marketing consents, got %+v", listBody.Data)
	}
}

func TestLegalDocumentsPublicAPI(t *testing.T) {
	app := newTestAppWithStore(memoryStore(t), "secret-token")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/legal/documents/personal_data", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body struct {
		Data struct {
			Version string `json:"version"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data.Version != "1.0" || body.Data.Content == "" {
		t.Fatalf("unexpected document: %+v", body.Data)
	}
}

func TestConsentVersionIsServerSide(t *testing.T) {
	app := newTestAppWithStore(memoryStore(t), "secret-token")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/consents", bytes.NewBufferString(`{
		"consent_type": "cookie_essential",
		"document_version": "99.0",
		"accepted_at": "2000-01-01T00:00:00Z"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var body struct {
		Data struct {
			DocumentVersion string `json:"document_version"`
			AcceptedAt      string `json:"accepted_at"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data.DocumentVersion != "1.0" {
		t.Fatalf("client must not set version, got %s", body.Data.DocumentVersion)
	}
	if body.Data.AcceptedAt == "2000-01-01T00:00:00Z" {
		t.Fatal("client must not set accepted_at")
	}
}

func memoryStore(t *testing.T) *memory.Store {
	t.Helper()
	return memory.NewStore()
}
