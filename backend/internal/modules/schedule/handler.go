package schedule

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the Schedule module.
type Handler struct {
	service Service
}

// NewHandler creates a new HTTP handler for schedule operations.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all schedule endpoints under the provided router group (e.g. /api/v1/schedule).
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Availability
	rg.GET("/availability", h.GetAvailability)
	rg.GET("/availability/:day", h.GetAvailabilityByDay)
	rg.PUT("/availability/:day", h.UpdateAvailabilityForDay)
	rg.DELETE("/availability/:day", h.DeleteAvailabilityForDay)

	// Sessions & occurrences
	rg.GET("/sessions", h.GetSessions)
	rg.POST("/sessions", h.CreateSession)
	rg.GET("/sessions/:id", h.GetSessionByID)
	rg.PUT("/sessions/:id", h.UpdateSession)
	rg.DELETE("/sessions/:id", h.ArchiveSession)

	// Exceptions
	rg.POST("/sessions/:id/exceptions", h.CreateException)
	rg.GET("/exceptions", h.ListExceptions)
	rg.DELETE("/sessions/:id/exceptions/:exceptionId", h.DeleteException)
}

// ==========================================
// AVAILABILITY ENDPOINTS
// ==========================================

// GetAvailability handles GET /api/v1/schedule/availability
func (h *Handler) GetAvailability(c *gin.Context) {
	avail, err := h.service.GetAvailability(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	if avail == nil {
		avail = make([]DayAvailabilityResponse, 0)
	}
	c.JSON(http.StatusOK, avail)
}

// GetAvailabilityByDay handles GET /api/v1/schedule/availability/:day
func (h *Handler) GetAvailabilityByDay(c *gin.Context) {
	day := c.Param("day")
	avail, err := h.service.GetAvailabilityByDay(c.Request.Context(), day)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, avail)
}

// UpdateAvailabilityForDay handles PUT /api/v1/schedule/availability/:day
func (h *Handler) UpdateAvailabilityForDay(c *gin.Context) {
	day := c.Param("day")
	var req UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido",
			},
		})
		return
	}

	avail, err := h.service.UpdateAvailabilityForDay(c.Request.Context(), day, req)
	if err != nil {
		var conflictErr *AvailabilityConflictError
		if errors.As(err, &conflictErr) {
			conflicts := conflictErr.Conflicts
			if conflicts == nil {
				conflicts = make([]SessionResponse, 0)
			}
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "AVAILABILITY_CONFLICT",
					"message": conflictErr.Message,
				},
				"conflicts": conflicts,
			})
			return
		}
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, avail)
}

// DeleteAvailabilityForDay handles DELETE /api/v1/schedule/availability/:day
func (h *Handler) DeleteAvailabilityForDay(c *gin.Context) {
	day := c.Param("day")
	if err := h.service.DeleteAvailabilityForDay(c.Request.Context(), day); err != nil {
		var conflictErr *AvailabilityConflictError
		if errors.As(err, &conflictErr) {
			conflicts := conflictErr.Conflicts
			if conflicts == nil {
				conflicts = make([]SessionResponse, 0)
			}
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "AVAILABILITY_CONFLICT",
					"message": conflictErr.Message,
				},
				"conflicts": conflicts,
			})
			return
		}
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disponibilidad del día eliminada exitosamente",
		"day":     day,
	})
}

// ==========================================
// SESSIONS ENDPOINTS
// ==========================================

// GetSessions handles GET /api/v1/schedule/sessions
// When 'from' and 'to' are provided, projects occurrences. If omitted, uses current week.
func (h *Handler) GetSessions(c *gin.Context) {
	filter := SessionsFilter{
		From:      c.Query("from"),
		To:        c.Query("to"),
		StudentID: c.Query("studentId"),
	}

	occurrences, err := h.service.GetOccurrences(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if occurrences == nil {
		occurrences = make([]SessionOccurrenceResponse, 0)
	}
	c.JSON(http.StatusOK, occurrences)
}

// CreateSession handles POST /api/v1/schedule/sessions
func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido o faltan campos obligatorios",
			},
		})
		return
	}

	session, err := h.service.CreateSession(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetSessionByID handles GET /api/v1/schedule/sessions/:id
func (h *Handler) GetSessionByID(c *gin.Context) {
	id := c.Param("id")
	session, err := h.service.GetSessionByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, session)
}

// UpdateSession handles PUT /api/v1/schedule/sessions/:id
func (h *Handler) UpdateSession(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido",
			},
		})
		return
	}

	session, err := h.service.UpdateSession(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, session)
}

// ArchiveSession handles DELETE /api/v1/schedule/sessions/:id
func (h *Handler) ArchiveSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.ArchiveSession(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesión finalizada / archivada exitosamente",
		"id":      id,
	})
}

// ==========================================
// EXCEPTIONS ENDPOINTS
// ==========================================

// CreateException handles POST /api/v1/schedule/sessions/:id/exceptions
func (h *Handler) CreateException(c *gin.Context) {
	sessionID := c.Param("id")
	var req CreateExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido o faltan campos obligatorios",
			},
		})
		return
	}

	exc, err := h.service.CreateException(c.Request.Context(), sessionID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, exc)
}

// ListExceptions handles GET /api/v1/schedule/exceptions
func (h *Handler) ListExceptions(c *gin.Context) {
	filter := ExceptionsFilter{
		From: c.Query("from"),
		To:   c.Query("to"),
	}

	list, err := h.service.ListExceptions(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if list == nil {
		list = make([]ExceptionResponse, 0)
	}
	c.JSON(http.StatusOK, list)
}

// DeleteException handles DELETE /api/v1/schedule/sessions/:id/exceptions/:exceptionId
func (h *Handler) DeleteException(c *gin.Context) {
	sessionID := c.Param("id")
	exceptionID := c.Param("exceptionId")

	if err := h.service.DeleteException(c.Request.Context(), sessionID, exceptionID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Excepción eliminada exitosamente; sesión recurrente restablecida",
		"sessionId":   sessionID,
		"exceptionId": exceptionID,
	})
}

// ==========================================
// ERROR MAPPING
// ==========================================

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrSessionNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "SESSION_NOT_FOUND",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrExceptionNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "EXCEPTION_NOT_FOUND",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrDuplicateException):
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "DUPLICATE_EXCEPTION",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrSessionOverlap):
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "SESSION_OVERLAP",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrSessionOutsideAvailability):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "OUTSIDE_AVAILABILITY",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrOutsideOperatingHours):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "OUTSIDE_OPERATING_HOURS",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrOverlappingIntervals):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "OVERLAPPING_INTERVALS",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrEndTimeBeforeStartTime),
		errors.Is(err, ErrSessionDurationTooShort),
		errors.Is(err, ErrInvalidTimeFormat),
		errors.Is(err, ErrInvalidDateFormat),
		errors.Is(err, ErrInvalidDateRange),
		errors.Is(err, ErrInvalidDayOfWeek),
		errors.Is(err, ErrInvalidExceptionType),
		errors.Is(err, ErrRescheduledMissingData),
		errors.Is(err, ErrDuplicateStudentInSession),
		errors.Is(err, ErrStudentNotFoundOrArchived):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrDatabaseUnavailable):
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "DATABASE_UNAVAILABLE",
				"message": err.Error(),
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Error interno del servidor",
			},
		})
	}
}
