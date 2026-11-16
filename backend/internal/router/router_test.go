package router

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"inverfit-backend/internal/config"
	"inverfit-backend/internal/handlers"
)

type mockPinger struct{}

func (m *mockPinger) Ping(ctx context.Context) error {
	return nil
}

func testConfig() *config.Config {
	return &config.Config{
		AppEnv:             "test",
		HTTPPort:           "8080",
		MongoDBURI:         "mongodb://localhost:27017",
		MongoDBDatabase:    "inverfit_test",
		CORSAllowedOrigins: []string{"http://localhost:5173"},
		LogLevel:           "info",
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRouterHealthLive(t *testing.T) {
	r := Setup(testConfig(), &mockPinger{}, nil, nil, nil, testLogger())

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	var resp handlers.LiveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al parsear JSON: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status esperado 'ok', obtenido '%s'", resp.Status)
	}
}

func TestRouterNotFound(t *testing.T) {
	r := Setup(testConfig(), &mockPinger{}, nil, nil, nil, testLogger())

	req := httptest.NewRequest(http.MethodGet, "/ruta/inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("código esperado 404, obtenido %d", w.Code)
	}

	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al parsear JSON: %v", err)
	}
	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("código de error esperado 'NOT_FOUND', obtenido '%s'", resp.Error.Code)
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	r := Setup(testConfig(), &mockPinger{}, nil, nil, nil, testLogger())

	// /health/live solo acepta GET, enviamos POST
	req := httptest.NewRequest(http.MethodPost, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("código esperado 405, obtenido %d", w.Code)
	}

	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al parsear JSON: %v", err)
	}
	if resp.Error.Code != "METHOD_NOT_ALLOWED" {
		t.Errorf("código de error esperado 'METHOD_NOT_ALLOWED', obtenido '%s'", resp.Error.Code)
	}
}
