package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCORSAllowedOrigin(t *testing.T) {
	r := gin.New()
	allowed := []string{"http://localhost:5173", "https://inverfit.pages.dev"}
	r.Use(CORS(allowed))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Test con origen permitido
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("código esperado 200, obtenido %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin esperado 'http://localhost:5173', obtenido '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSDisallowedOrigin(t *testing.T) {
	r := gin.New()
	allowed := []string{"http://localhost:5173"}
	r.Use(CORS(allowed))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Test con origen no permitido
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://sitio-malicioso.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("origen no permitido no debe tener Access-Control-Allow-Origin, obtenido: %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSPreflightOptions(t *testing.T) {
	r := gin.New()
	allowed := []string{"http://localhost:5173"}
	r.Use(CORS(allowed))

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("preflight OPTIONS debe retornar 204 No Content, obtenido %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Errorf("se esperaban métodos permitidos en cabecera Access-Control-Allow-Methods")
	}
}

func TestRecoveryFromPanic(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))

	r := gin.New()
	r.Use(Recovery(logger))

	r.GET("/panic", func(c *gin.Context) {
		panic("error inesperado de prueba")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("código esperado 500, obtenido %d", w.Code)
	}

	var resp map[string]map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al parsear JSON de respuesta: %v", err)
	}

	if resp["error"]["code"] != "INTERNAL_SERVER_ERROR" {
		t.Errorf("código de error esperado 'INTERNAL_SERVER_ERROR', obtenido '%s'", resp["error"]["code"])
	}
	if resp["error"]["message"] != "Ha ocurrido un error interno en el servidor" {
		t.Errorf("mensaje inesperado: %s", resp["error"]["message"])
	}

	// Verificar que el log del servidor sí registró el mensaje del panic
	if !bytes.Contains(logBuf.Bytes(), []byte("error inesperado de prueba")) {
		t.Errorf("el log del servidor debió registrar el detalle del panic")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))

	r := gin.New()
	r.Use(Logging(logger))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/health?token=secret123&api_key=private456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	logOutput := logBuf.String()
	if !bytes.Contains([]byte(logOutput), []byte("/health")) {
		t.Errorf("el log debe contener la ruta /health, obtenido: %s", logOutput)
	}
	if bytes.Contains([]byte(logOutput), []byte("secret123")) || bytes.Contains([]byte(logOutput), []byte("private456")) {
		t.Errorf("el log no debe contener parámetros de consulta sensibles (query strings), obtenido: %s", logOutput)
	}
}
