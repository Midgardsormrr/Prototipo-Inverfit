package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockPinger simula la interfaz Pinger para pruebas unitarias sin requerir MongoDB en ejecución.
type mockPinger struct {
	pingErr error
}

func (m *mockPinger) Ping(ctx context.Context) error {
	return m.pingErr
}

func TestHealthLive(t *testing.T) {
	h := NewHealthHandler(&mockPinger{})

	r := gin.New()
	r.GET("/health/live", h.Live)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	var resp LiveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al deserializar JSON: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("status esperado 'ok', obtenido '%s'", resp.Status)
	}
	if resp.Service != "inverfit-api" {
		t.Errorf("service esperado 'inverfit-api', obtenido '%s'", resp.Service)
	}
}

func TestHealthReadyConnected(t *testing.T) {
	h := NewHealthHandler(&mockPinger{pingErr: nil})

	r := gin.New()
	r.GET("/health/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	var resp ReadyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al deserializar JSON: %v", err)
	}

	if resp.Status != "ready" {
		t.Errorf("status esperado 'ready', obtenido '%s'", resp.Status)
	}
	if resp.Service != "inverfit-api" {
		t.Errorf("service esperado 'inverfit-api', obtenido '%s'", resp.Service)
	}
	if resp.Dependencies["mongodb"] != "connected" {
		t.Errorf("dependencia mongodb esperada 'connected', obtenido '%s'", resp.Dependencies["mongodb"])
	}
}

func TestHealthReadyDisconnected(t *testing.T) {
	h := NewHealthHandler(&mockPinger{pingErr: errors.New("timeout de conexión al socket")})

	r := gin.New()
	r.GET("/health/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("código esperado 503, obtenido %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al deserializar JSON de error: %v", err)
	}

	if resp.Error.Code != "SERVICE_UNAVAILABLE" {
		t.Errorf("código de error esperado 'SERVICE_UNAVAILABLE', obtenido '%s'", resp.Error.Code)
	}
	if resp.Error.Message != "El servicio no se encuentra disponible" {
		t.Errorf("mensaje de error esperado 'El servicio no se encuentra disponible', obtenido '%s'", resp.Error.Message)
	}
}

func TestHealthReadyNilPinger(t *testing.T) {
	h := NewHealthHandler(nil)

	r := gin.New()
	r.GET("/health/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("código esperado 503, obtenido %d", w.Code)
	}
}
