package tests

import (
	"burrowfs/api/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type statusResponse struct {
	Status string
}

func TestStatusRoute(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler := routes.StatusRoutes()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response statusResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response.Status)
	}
}
