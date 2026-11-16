package students

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorDetail estructura el detalle de error estandarizado para la API.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse estructura la respuesta uniforme de error.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Handler gestiona los endpoints HTTP del módulo de alumnos.
type Handler struct {
	service Service
}

// NewHandler inicializa una nueva instancia de Handler con su servicio inyectado.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registra los endpoints del módulo de alumnos bajo el grupo recibido (ej: /api/v1/students).
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.POST("", h.Create)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

// List maneja GET /api/v1/students.
func (h *Handler) List(c *gin.Context) {
	filter := StudentListFilter{
		Name:      c.Query("name"),
		Status:    c.Query("status"),
		Frequency: c.Query("frequency"),
	}

	result, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByID maneja GET /api/v1/students/:id.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	student, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, student)
}

// Create maneja POST /api/v1/students.
func (h *Handler) Create(c *gin.Context) {
	var req CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_JSON",
				Message: "Formato de cuerpo JSON inválido o malformado",
			},
		})
		return
	}

	created, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

// Update maneja PUT /api/v1/students/:id.
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_JSON",
				Message: "Formato de cuerpo JSON inválido o malformado",
			},
		})
		return
	}

	updated, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete maneja DELETE /api/v1/students/:id (eliminación lógica / archivado).
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alumno archivado exitosamente",
	})
}

// respondError mapea los errores de dominio de la capa service a respuestas HTTP con códigos de estado pertinentes.
func (h *Handler) respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrStudentNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Code:    "STUDENT_NOT_FOUND",
				Message: "El alumno solicitado no existe o se encuentra archivado",
			},
		})
	case errors.Is(err, ErrInvalidID):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_ID",
				Message: "El identificador de alumno proporcionado no tiene un formato válido",
			},
		})
	case errors.Is(err, ErrEmptyName):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "EMPTY_NAME",
				Message: err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidFrequency):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_FREQUENCY",
				Message: err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_STATUS",
				Message: err.Error(),
			},
		})
	case errors.Is(err, ErrDatabaseUnavailable):
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "Base de datos no disponible para completar la solicitud",
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Ocurrió un error inesperado al procesar la solicitud",
			},
		})
	}
}
