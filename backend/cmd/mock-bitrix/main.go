// mock-bitrix is a minimal in-memory Bitrix24 REST mock for local integration testing.
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
	addr := envString("MOCK_BITRIX_ADDR", ":8091")
	store := newStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/debug/deals/", store.handleDebugDealStage)
	mux.HandleFunc("/rest/", store.handleREST)

	log.Printf("mock-bitrix listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

type store struct {
	mu sync.Mutex

	nextContactID int
	nextDealID    int
	nextProductID int
	nextCommentID int

	contacts []map[string]any
	deals    []map[string]any
	products []map[string]any
}

func newStore() *store {
	return &store{
		nextContactID: 100,
		nextDealID:    1000,
		nextProductID: 500,
		nextCommentID: 9000,
	}
}

func (s *store) handleREST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	method := restMethod(r.URL.Path)
	if method == "" {
		http.Error(w, "unknown rest path", http.StatusNotFound)
		return
	}

	var params map[string]any
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch method {
	case "crm.duplicate.findbycomm":
		writeResult(w, s.findDuplicate(params))
	case "crm.contact.add":
		writeResult(w, s.addContact(params))
	case "crm.deal.list":
		writeResult(w, s.listDeals(params))
	case "crm.deal.add":
		writeResult(w, s.addDeal(params))
	case "crm.deal.update":
		writeResult(w, s.updateDeal(params))
	case "crm.deal.get":
		writeResult(w, s.getDeal(params))
	case "crm.product.list":
		writeResult(w, s.listProducts(params))
	case "crm.product.add":
		writeResult(w, s.addProduct(params))
	case "crm.product.update":
		writeResult(w, s.updateProduct(params))
	case "crm.deal.productrows.set":
		writeResult(w, true)
	case "crm.timeline.comment.add":
		writeResult(w, s.addComment())
	default:
		http.Error(w, "unknown method: "+method, http.StatusNotFound)
	}
}

func (s *store) handleDebugDealStage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/debug/deals/")
	id = strings.TrimSuffix(id, "/stage")
	if id == "" {
		http.Error(w, "deal id required", http.StatusBadRequest)
		return
	}

	var body struct {
		StageID string `json:"stage_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.StageID == "" {
		http.Error(w, "stage_id required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, deal := range s.deals {
		if fmt.Sprintf("%v", deal["ID"]) == id {
			s.deals[i]["STAGE_ID"] = body.StageID
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "deal not found", http.StatusNotFound)
}

func (s *store) findDuplicate(params map[string]any) map[string][]any {
	values := stringSlice(params["values"])
	for _, contact := range s.contacts {
		phones := contactPhones(contact)
		for _, phone := range values {
			for _, existing := range phones {
				if existing == phone {
					return map[string][]any{"CONTACT": {contact["ID"]}}
				}
			}
		}
	}
	return map[string][]any{"CONTACT": {}}
}

func (s *store) addContact(params map[string]any) int {
	s.nextContactID++
	fields, _ := params["fields"].(map[string]any)
	contact := map[string]any{"ID": s.nextContactID}
	for k, v := range fields {
		contact[k] = v
	}
	s.contacts = append(s.contacts, contact)
	return s.nextContactID
}

func (s *store) listDeals(params map[string]any) []map[string]any {
	filter, _ := params["filter"].(map[string]any)
	var items []map[string]any
	for _, deal := range s.deals {
		if matchFilter(deal, filter) {
			items = append(items, map[string]any{"ID": deal["ID"]})
		}
	}
	if items == nil {
		return []map[string]any{}
	}
	return items
}

func (s *store) addDeal(params map[string]any) int {
	s.nextDealID++
	fields, _ := params["fields"].(map[string]any)
	deal := map[string]any{"ID": s.nextDealID}
	for k, v := range fields {
		deal[k] = v
	}
	if _, ok := deal["STAGE_ID"]; !ok {
		deal["STAGE_ID"] = "NEW"
	}
	s.deals = append(s.deals, deal)
	return s.nextDealID
}

func (s *store) updateDeal(params map[string]any) bool {
	id := fmt.Sprintf("%v", params["id"])
	fields, _ := params["fields"].(map[string]any)
	for i, deal := range s.deals {
		if fmt.Sprintf("%v", deal["ID"]) == id {
			for k, v := range fields {
				s.deals[i][k] = v
			}
			return true
		}
	}
	return false
}

func (s *store) getDeal(params map[string]any) map[string]any {
	id := fmt.Sprintf("%v", params["id"])
	for _, deal := range s.deals {
		if fmt.Sprintf("%v", deal["ID"]) == id {
			return deal
		}
	}
	return map[string]any{}
}

func (s *store) listProducts(params map[string]any) []map[string]any {
	filter, _ := params["filter"].(map[string]any)
	var items []map[string]any
	for _, product := range s.products {
		if matchFilter(product, filter) {
			items = append(items, map[string]any{"ID": product["ID"]})
		}
	}
	if items == nil {
		return []map[string]any{}
	}
	return items
}

func (s *store) addProduct(params map[string]any) int {
	s.nextProductID++
	fields, _ := params["fields"].(map[string]any)
	product := map[string]any{"ID": s.nextProductID}
	for k, v := range fields {
		product[k] = v
	}
	s.products = append(s.products, product)
	return s.nextProductID
}

func (s *store) updateProduct(params map[string]any) bool {
	id := fmt.Sprintf("%v", params["id"])
	fields, _ := params["fields"].(map[string]any)
	for i, product := range s.products {
		if fmt.Sprintf("%v", product["ID"]) == id {
			for k, v := range fields {
				s.products[i][k] = v
			}
			return true
		}
	}
	return false
}

func (s *store) addComment() int {
	s.nextCommentID++
	return s.nextCommentID
}

func matchFilter(entity map[string]any, filter map[string]any) bool {
	if len(filter) == 0 {
		return true
	}
	for key, expected := range filter {
		actual, ok := entity[key]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
			return false
		}
	}
	return true
}

func contactPhones(contact map[string]any) []string {
	raw, ok := contact["PHONE"].([]any)
	if !ok {
		return nil
	}
	var phones []string
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if v, ok := m["VALUE"].(string); ok {
			phones = append(phones, v)
		}
	}
	return phones
}

func stringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, fmt.Sprintf("%v", item))
	}
	return out
}

func writeResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"result": result})
}

func envString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func restMethod(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.Contains(parts[i], ".") {
			return parts[i]
		}
	}
	return ""
}
