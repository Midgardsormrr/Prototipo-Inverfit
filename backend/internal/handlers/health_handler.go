package handlers

import (
	"context"
	"net/http"
	"time"

	"inverfit-backend/internal/database"

	"github.com/gin-gonic/gin"
)

// Constante de identificación del servicio
const ServiceName = "inverfit-api"

// HealthHandler administra los endpoints de comprobación de salud del proceso y de sus dependencias.
type HealthHandler struct {
	pinger database.Pinger
}

// NewHealthHandler crea una nueva instancia de HealthHandler con la interfaz de comprobación de base de datos.
func NewHealthHandler(pinger database.Pinger) *HealthHandler {
	return &HealthHandler{
		pinger: pinger,
	}
}

// LiveResponse estructura la respuesta del endpoint /health/live.
type LiveResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// ReadyResponse estructura la respuesta exitosa del endpoint /health/ready.
type ReadyResponse struct {
	Status       string            `json:"status"`
	Service      string            `json:"service"`
	Dependencies map[string]string `json:"dependencies"`
}

// ErrorDetail estructura el detalle de un error estandarizado de la API.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse estructura la respuesta de error uniforme de la API.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Live maneja GET /health/live.
// Comprueba que el proceso de la API está en ejecución y puede recibir peticiones HTTP.
// No tiene dependencias externas (no consulta a MongoDB).
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, LiveResponse{
		Status:  "ok",
		Service: ServiceName,
	})
}

// Ready maneja GET /health/ready.
// Comprueba que el backend está preparado para procesar tráfico verificando la conectividad efectiva con MongoDB.
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.pinger == nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "El servicio no se encuentra disponible",
			},
		})
		return
	}

	// Timeout breve para la verificación de salud de la base de datos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.pinger.Ping(ctx); err != nil {
		// Respondemos 503 sin exponer detalles internos ni credenciales
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "El servicio no se encuentra disponible",
			},
		})
		return
	}

	c.JSON(http.StatusOK, ReadyResponse{
		Status:  "ready",
		Service: ServiceName,
		Dependencies: map[string]string{
			"mongodb": "connected",
		},
	})
}
