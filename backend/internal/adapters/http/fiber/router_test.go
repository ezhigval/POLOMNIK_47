package fiberhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"polomnik/internal/adapters/integration/noop"
	notificationnoop "polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/application"
	"polomnik/internal/config"
	"polomnik/internal/domain"
)

func TestHealthRoutes(t *testing.T) {
	app := newTestApp(config.Config{AppEnv: "test", HTTPAddr: ":0"})

	for _, path := range []string{"/health", "/api/v1/health", "/health/ready", "/api/v1/health/ready"} {
		req := httptest.NewRequest("GET", path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request %s failed: %v", path, err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("expected %s to return 200, got %d", path, resp.StatusCode)
		}

		var body struct {
			Data struct {
				Status  string `json:"status"`
				Service string `json:"service"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("decode %s response: %v", path, err)
		}
		if body.Data.Status != "ok" {
			t.Fatalf("expected %s status ok, got %q", path, body.Data.Status)
		}
	}
}

func TestNotFoundUsesErrorEnvelope(t *testing.T) {
	app := newTestApp(config.Config{AppEnv: "test", HTTPAddr: ":0"})

	req := httptest.NewRequest("GET", "/missing", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %q", body.Error.Code)
	}
}

func TestCreateBookingAndListTours(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")
	start := time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02")
	end := time.Now().UTC().AddDate(0, 0, 34).Format("2006-01-02")

	createTourReq := httptest.NewRequest(http.MethodPost, "/api/v1/management/tours", bytes.NewBufferString(`{
		"slug": "pilgrimage",
		"title": "Pilgrimage Tour",
		"description": "Test",
		"price": 15000,
		"currency": "RUB",
		"date_start": "`+start+`",
		"date_end": "`+end+`",
		"slots_total": 10,
		"slots_left": 10,
		"location": "Moscow",
		"images": [],
		"is_active": true,
		"is_hot": true,
		"overbooking_enabled": false
	}`))
	createTourReq.Header.Set("Content-Type", "application/json")
	createTourReq.Header.Set("X-Admin-Token", "secret-token")

	createTourResp, err := app.Test(createTourReq)
	if err != nil {
		t.Fatalf("create tour request failed: %v", err)
	}
	if createTourResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createTourResp.StatusCode)
	}

	var createdTour struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(createTourResp.Body).Decode(&createdTour); err != nil {
		t.Fatalf("decode create tour response: %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tours", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list tours request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}

	bookingReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(`{
		"tour_id": "`+createdTour.Data.ID+`",
		"name": "Иван Иванов",
		"phone": "+79999999999",
		"email": "mail@test.com",
		"people_count": 2,
		"comment": "Please call me"
	}`))
	bookingReq.Header.Set("Content-Type", "application/json")

	bookingResp, err := app.Test(bookingReq)
	if err != nil {
		t.Fatalf("create booking request failed: %v", err)
	}
	if bookingResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", bookingResp.StatusCode)
	}

	var bookingBody struct {
		Data struct {
			BookingID         string `json:"booking_id"`
			BookingStatus     string `json:"booking_status"`
			TotalPrice        int    `json:"total_price"`
			IntegrationStatus string `json:"integration_status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(bookingResp.Body).Decode(&bookingBody); err != nil {
		t.Fatalf("decode booking response: %v", err)
	}
	if bookingBody.Data.BookingStatus != "NEW" {
		t.Fatalf("expected NEW status, got %q", bookingBody.Data.BookingStatus)
	}
	if bookingBody.Data.TotalPrice != 30000 {
		t.Fatalf("expected total price 30000, got %d", bookingBody.Data.TotalPrice)
	}
	if bookingBody.Data.IntegrationStatus != "not_configured" {
		t.Fatalf("expected not_configured, got %q", bookingBody.Data.IntegrationStatus)
	}
}

func TestManagementListIntegrationReferences(t *testing.T) {
	app := newTestApp(config.Config{AppEnv: "test", HTTPAddr: ":0", AdminToken: "secret"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/management/integration-references", nil)
	req.Header.Set("X-Admin-Token", "secret")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("list integration references request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestManagementRequiresAdminToken(t *testing.T) {
	app := newTestAppWithStore(memory.NewStore(), "secret-token")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/management/tours", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestCreateBookingValidationError(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")

	tour := testTourForHandler(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(`{
		"tour_id": "`+tour.ID.String()+`",
		"name": "Test",
		"phone": "+79999999999",
		"people_count": 0
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestListToursMetaUsesSnakeCase(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")

	tour := testTourForHandler(t, store)
	_ = tour

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tours", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Meta struct {
			Page    int  `json:"page"`
			HasNext bool `json:"has_next"`
		} `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Meta.Page != 1 {
		t.Fatalf("expected page=1, got %d", body.Meta.Page)
	}
}

func testTourForHandler(t *testing.T, store *memory.Store) domain.Tour {
	t.Helper()

	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "handler-test-tour",
		Title:      "Handler Test Tour",
		Price:      15000,
		Currency:   "RUB",
		DateStart:  time.Now().UTC().AddDate(0, 0, 30),
		DateEnd:    time.Now().UTC().AddDate(0, 0, 34),
		SlotsTotal: 10,
		SlotsLeft:  10,
		Location:   "Moscow",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if _, err := store.CreateTour(context.Background(), tour); err != nil {
		t.Fatalf("store tour: %v", err)
	}

	return tour
}

func TestOAuthRequiresInternalSecret(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/oauth", bytes.NewBufferString(`{
		"provider": "google",
		"subject": "123",
		"email": "attacker@example.com",
		"name": "Attacker"
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without internal secret, got %d", resp.StatusCode)
	}
}

func TestCreateAndListNews(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/management/news", bytes.NewBufferString(`{
		"slug": "tikhvin-path",
		"title": "Тихвинский путь",
		"excerpt": "Анонс",
		"body": "Первый абзац.\n\nВторой абзац.",
		"image_url": "/images/hero/tikhvin-monastery.webp",
		"published_at": "2026-08-21",
		"is_published": true,
		"sort_order": 0
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-Admin-Token", "secret-token")

	createResp, err := app.Test(createReq)
	if err != nil {
		t.Fatalf("create news request failed: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.StatusCode)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/news", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list news request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}

	var body struct {
		Data []struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
		} `json:"data"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&body); err != nil {
		t.Fatalf("decode list news: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Slug != "tikhvin-path" {
		t.Fatalf("expected one published article, got %+v", body.Data)
	}
}

func TestSetCompanyReplyOnReview(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "secret-token")

	createTourReq := httptest.NewRequest(http.MethodPost, "/api/v1/management/tours", bytes.NewBufferString(`{
		"slug": "optina",
		"title": "Оптина пустынь",
		"description": "Test",
		"price": 15000,
		"currency": "RUB",
		"date_start": "2026-08-01",
		"date_end": "2026-08-05",
		"slots_total": 10,
		"slots_left": 10,
		"location": "Козельск",
		"images": [],
		"is_active": true,
		"is_hot": false,
		"overbooking_enabled": false
	}`))
	createTourReq.Header.Set("Content-Type", "application/json")
	createTourReq.Header.Set("X-Admin-Token", "secret-token")

	createTourResp, err := app.Test(createTourReq)
	if err != nil {
		t.Fatalf("create tour request failed: %v", err)
	}
	if createTourResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createTourResp.StatusCode)
	}

	var createdTour struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(createTourResp.Body).Decode(&createdTour); err != nil {
		t.Fatalf("decode created tour: %v", err)
	}

	createReviewBody := fmt.Sprintf(`{
		"tour_id": %q,
		"client_name": "Мария К.",
		"rating": 5,
		"text": "Поездка оставила глубокое впечатление.",
		"is_approved": true
	}`, createdTour.Data.ID)
	createReviewReq := httptest.NewRequest(http.MethodPost, "/api/v1/management/reviews", bytes.NewBufferString(createReviewBody))
	createReviewReq.Header.Set("Content-Type", "application/json")
	createReviewReq.Header.Set("X-Admin-Token", "secret-token")

	createReviewResp, err := app.Test(createReviewReq)
	if err != nil {
		t.Fatalf("create review request failed: %v", err)
	}
	if createReviewResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createReviewResp.StatusCode)
	}

	var createdReview struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(createReviewResp.Body).Decode(&createdReview); err != nil {
		t.Fatalf("decode created review: %v", err)
	}

	replyReq := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/management/reviews/"+createdReview.Data.ID+"/reply",
		bytes.NewBufferString(`{"company_reply":"Благодарим за добрые слова."}`),
	)
	replyReq.Header.Set("Content-Type", "application/json")
	replyReq.Header.Set("X-Admin-Token", "secret-token")

	replyResp, err := app.Test(replyReq)
	if err != nil {
		t.Fatalf("set reply request failed: %v", err)
	}
	if replyResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", replyResp.StatusCode)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/reviews", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list reviews request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}

	var body struct {
		Data []struct {
			CompanyReply string `json:"company_reply"`
		} `json:"data"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&body); err != nil {
		t.Fatalf("decode list reviews: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].CompanyReply != "Благодарим за добрые слова." {
		t.Fatalf("expected public review with company reply, got %+v", body.Data)
	}
}

func newTestApp(cfg config.Config) *fiber.App {
	return newTestAppWithStore(memory.NewStore(), cfg.AdminToken)
}

func newTestAppWithStore(store *memory.Store, adminToken string) *fiber.App {
	bookingService := application.NewBookingService(
		store,
		store,
		noop.NewCRMAdapter(),
		noop.NewAccountingAdapter(),
		notificationnoop.New(),
		store,
	)
	return NewRouter(config.Config{
		AppEnv:            "test",
		HTTPAddr:          ":0",
		AdminToken:        adminToken,
		InternalAPISecret: config.DefaultInternalAPISecret,
		JWTSecret:         config.DefaultJWTSecret,
	}, discardLogger(), Services{
		Tours:        application.NewTourService(store, nil, noop.NewCRMAdapter()),
		Bookings:     bookingService,
		Reviews:      application.NewReviewService(store, store, noop.NewCRMAdapter()),
		Integrations: application.NewIntegrationService(store, store),
		Webhooks: application.NewWebhookService(
			bookingService,
			noop.NewCRMInboundAdapter(),
			"",
			false,
		),
		Auth: application.NewAuthService(store, store, nil, nil, application.SocialAuthConfig{}, config.DefaultJWTSecret, 24*time.Hour, "http://localhost:3000"),
		Support:  application.NewSupportService(store, notificationnoop.New()),
		CMS:      application.NewCMSService(store),
		News:     application.NewNewsService(store),
		Telegram: application.NewTelegramServiceFromRepos(store, store, nil, ""),
		Notifications: application.NewNotificationSettingsService(store, store, store, false, false),
		SiteSettings:  application.NewSiteSettingsService(store, domain.SiteSettings{}),
		AdminRoles:    application.NewAdminRoleService(store, store, adminToken, config.DefaultJWTSecret),
	}, HealthDeps{})
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestManagementListCMSPagesEmpty(t *testing.T) {
	app := newTestAppWithStore(memory.NewStore(), "admin-token")

	req := httptest.NewRequest("GET", "/api/v1/management/cms/pages", nil)
	req.Header.Set("X-Admin-Token", "admin-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d %s", resp.StatusCode, body)
	}

	var body struct {
		Data []struct {
			ID    string `json:"id"`
			Slug  string `json:"slug"`
			Title string `json:"title"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data == nil {
		t.Fatal("cms pages data must be an array, not null")
	}
}

func TestPublicCMSPageMissingIs404(t *testing.T) {
	app := newTestApp(config.Config{AppEnv: "test", HTTPAddr: ":0"})

	req := httptest.NewRequest("GET", "/api/v1/pages/home", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 404 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 404 for unpublished/missing home, got %d %s", resp.StatusCode, body)
	}
}

func TestTelegramSettingsAndWebhook(t *testing.T) {
	store := memory.NewStore()
	app := newTestAppWithStore(store, "admin-token")

	patch := httptest.NewRequest("PATCH", "/api/v1/management/telegram-settings", bytes.NewBufferString(`{"booking_usernames":"@EzhigVal","support_usernames":"other_user"}`))
	patch.Header.Set("Content-Type", "application/json")
	patch.Header.Set("X-Admin-Token", "admin-token")
	patchResp, err := app.Test(patch)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if patchResp.StatusCode != 200 {
		body, _ := io.ReadAll(patchResp.Body)
		t.Fatalf("expected patch 200, got %d %s", patchResp.StatusCode, body)
	}

	invalid := httptest.NewRequest("PATCH", "/api/v1/management/telegram-settings", bytes.NewBufferString(`{"booking_usernames":"ab"}`))
	invalid.Header.Set("Content-Type", "application/json")
	invalid.Header.Set("X-Admin-Token", "admin-token")
	invalidResp, err := app.Test(invalid)
	if err != nil {
		t.Fatalf("invalid: %v", err)
	}
	if invalidResp.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", invalidResp.StatusCode)
	}

	webhook := httptest.NewRequest("POST", "/api/v1/webhooks/telegram", bytes.NewBufferString(`{"message":{"from":{"username":"EzhigVal"},"chat":{"id":111},"text":"/start"}}`))
	webhook.Header.Set("Content-Type", "application/json")
	webhook.Header.Set("X-Telegram-Bot-Api-Secret-Token", application.TelegramWebhookSecret(config.DefaultInternalAPISecret))
	webhookResp, err := app.Test(webhook)
	if err != nil {
		t.Fatalf("webhook: %v", err)
	}
	if webhookResp.StatusCode != 200 {
		body, _ := io.ReadAll(webhookResp.Body)
		t.Fatalf("expected webhook 200, got %d %s", webhookResp.StatusCode, body)
	}

	getReq := httptest.NewRequest("GET", "/api/v1/management/telegram-settings", nil)
	getReq.Header.Set("X-Admin-Token", "admin-token")
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if getResp.StatusCode != 200 {
		t.Fatalf("expected get 200, got %d", getResp.StatusCode)
	}
	var body struct {
		Data struct {
			Recipients []struct {
				Username  string `json:"username"`
				ChatBound bool   `json:"chat_bound"`
			} `json:"recipients"`
		} `json:"data"`
	}
	if err := json.NewDecoder(getResp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	bound := false
	for _, item := range body.Data.Recipients {
		if item.Username == "ezhigval" && item.ChatBound {
			bound = true
		}
	}
	if !bound {
		t.Fatalf("expected ezhigval to be bound after /start, got %+v", body.Data.Recipients)
	}
}
