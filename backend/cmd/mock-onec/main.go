// mock-onec is a minimal 1C HTTP exchange mock (see docs/ONEC_INTEGRATOR_TZ.md).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	addr := envString("MOCK_ONEC_ADDR", ":8092")
	basePath := strings.TrimRight(envString("MOCK_ONEC_BASE_PATH", "/accounting"), "/")
	username := os.Getenv("MOCK_ONEC_USERNAME")
	password := os.Getenv("MOCK_ONEC_PASSWORD")

	store := &docStore{
		bookingPath:       basePath + "/hs/polomnik/booking",
		counterpartyPath:  basePath + "/hs/polomnik/counterparty",
		nextBookingDoc:    1,
		nextCounterparty:  1,
		bookings:          make(map[string]exchangeResponse),
		counterparties:    make(map[string]exchangeResponse),
		phones:            make(map[string]string),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc(store.bookingPath, authWrap(username, password, store.handleBooking))
	mux.HandleFunc(store.counterpartyPath, authWrap(username, password, store.handleCounterparty))
	mux.HandleFunc(basePath+"/debug/documents", store.handleDebug)

	log.Printf("mock-onec listening on %s (base %s)", addr, basePath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

type exchangeResponse struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type bookingPayload struct {
	BookingID   string `json:"booking_id"`
	TourID      string `json:"tour_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	PeopleCount int    `json:"people_count"`
	TotalPrice  int    `json:"total_price"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
	Source      string `json:"source"`
	Overbooked  bool   `json:"overbooked"`
}

type counterpartyPayload struct {
	BookingID string `json:"booking_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type docStore struct {
	mu sync.Mutex

	bookingPath      string
	counterpartyPath string
	nextBookingDoc   int
	nextCounterparty int

	bookings       map[string]exchangeResponse
	counterparties map[string]exchangeResponse
	phones         map[string]string
}

func (s *docStore) handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload bookingPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.BookingID == "" {
		http.Error(w, "invalid booking payload", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.bookings[payload.BookingID]; ok {
		writeJSON(w, existing)
		return
	}

	s.nextBookingDoc++
	resp := exchangeResponse{
		ExternalID: formatDoc("DOC", s.nextBookingDoc),
		Status:     "ok",
		Message:    "booking exported",
	}
	s.bookings[payload.BookingID] = resp
	writeJSON(w, resp)
}

func (s *docStore) handleCounterparty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload counterpartyPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.BookingID == "" {
		http.Error(w, "invalid counterparty payload", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if payload.Phone != "" {
		if existingID, ok := s.phones[payload.Phone]; ok {
			resp := exchangeResponse{ExternalID: existingID, Status: "ok", Message: "counterparty reused"}
			s.counterparties[payload.BookingID] = resp
			writeJSON(w, resp)
			return
		}
	}

	if existing, ok := s.counterparties[payload.BookingID]; ok {
		writeJSON(w, existing)
		return
	}

	s.nextCounterparty++
	resp := exchangeResponse{
		ExternalID: formatDoc("CP", s.nextCounterparty),
		Status:     "ok",
		Message:    "counterparty created",
	}
	s.counterparties[payload.BookingID] = resp
	if payload.Phone != "" {
		s.phones[payload.Phone] = resp.ExternalID
	}
	writeJSON(w, resp)
}

func (s *docStore) handleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	writeJSON(w, map[string]any{
		"bookings":       s.bookings,
		"counterparties": s.counterparties,
	})
}

func authWrap(username, password string, next http.HandlerFunc) http.HandlerFunc {
	if username == "" {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="mock-onec"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func formatDoc(prefix string, n int) string {
	return fmt.Sprintf("%s-%06d", prefix, n)
}

func envString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
