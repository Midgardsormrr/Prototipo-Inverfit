package payments

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the payments module.
type Handler struct {
	service Service
}

// NewHandler creates a new payments Handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes mounts all payment endpoints onto the specified Gin router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListPayments)
	rg.POST("", h.RegisterPayment)
	rg.GET("/summary", h.GetFinancialSummary)
	rg.GET("/export", h.ExportCSV)
	rg.GET("/student/:studentId", h.GetStudentPayments)
	rg.GET("/:id", h.GetPaymentByID)
	rg.PUT("/:id", h.UpdatePayment)
	rg.DELETE("/:id", h.ArchivePayment)
	rg.POST("/:id/simulate-reminder", h.SimulateReminder)
	rg.POST("/:id/simulate-receipt", h.SimulateReceipt)
}

// ListPayments returns active payments with optional filters.
func (h *Handler) ListPayments(c *gin.Context) {
	filter := PaymentFilter{
		StudentID: c.Query("studentId"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		From:      c.Query("from"),
		To:        c.Query("to"),
	}

	payments, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// RegisterPayment handles registering a new payment or confirming a pending one.
func (h *Handler) RegisterPayment(c *gin.Context) {
	var req RegisterPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido o faltan campos obligatorios",
			},
		})
		return
	}

	payment, err := h.service.RegisterPayment(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetFinancialSummary returns aggregated KPIs calculated from active payments.
func (h *Handler) GetFinancialSummary(c *gin.Context) {
	summary, err := h.service.GetFinancialSummary(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetStudentPayments returns all active payments for a student.
func (h *Handler) GetStudentPayments(c *gin.Context) {
	studentID := c.Param("studentId")
	payments, err := h.service.GetByStudentID(c.Request.Context(), studentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentByID returns details of a single payment by ID.
func (h *Handler) GetPaymentByID(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// UpdatePayment partially modifies an existing payment.
func (h *Handler) UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_JSON",
				"message": "El cuerpo de la petición no contiene un JSON válido",
			},
		})
		return
	}

	payment, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ArchivePayment soft-deletes a payment record.
func (h *Handler) ArchivePayment(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Archive(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registro de pago archivado exitosamente",
		"id":      id,
	})
}

// SimulateReminder updates the payment reminder count and timestamp.
func (h *Handler) SimulateReminder(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.SimulateReminder(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// SimulateReceipt marks the receipt as sent with timestamp.
func (h *Handler) SimulateReceipt(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.SimulateReceipt(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ExportCSV streams a UTF-8 CSV file containing payments, protected against Formula Injection.
func (h *Handler) ExportCSV(c *gin.Context) {
	filter := PaymentFilter{
		StudentID: c.Query("studentId"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		From:      c.Query("from"),
		To:        c.Query("to"),
	}

	payments, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	todayStr := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("inverfit_pagos_%s.csv", todayStr)

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	// Write UTF-8 BOM so Excel opens it with proper UTF-8 encoding
	_, _ = c.Writer.Write([]byte("\xef\xbb\xbf"))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header row
	headers := []string{
		"ID",
		"Alumno",
		"Plan",
		"Frecuencia",
		"Vencimiento",
		"Monto",
		"Estado",
		"Medio de Pago",
		"Fecha Pago",
		"Comprobante Enviado",
		"Notas",
	}
	if err := writer.Write(headers); err != nil {
		return
	}

	// Write data rows
	for _, p := range payments {
		amountStr := "Por definir"
		if p.Amount != nil {
			amountStr = strconv.FormatInt(*p.Amount, 10)
		}

		paymentDateStr := "Pendiente"
		if p.PaymentDate != nil && *p.PaymentDate != "" {
			paymentDateStr = *p.PaymentDate
		}

		receiptSentStr := "No"
		if p.ReceiptSent {
			receiptSentStr = "Sí"
		}

		row := []string{
			SanitizeCSVCell(p.ID),
			SanitizeCSVCell(p.StudentName),
			SanitizeCSVCell(p.PlanName),
			SanitizeCSVCell(p.Frequency),
			SanitizeCSVCell(p.DueDate),
			amountStr,
			SanitizeCSVCell(p.Status),
			SanitizeCSVCell(p.PaymentMethod),
			SanitizeCSVCell(paymentDateStr),
			receiptSentStr,
			SanitizeCSVCell(p.Notes),
		}

		if err := writer.Write(row); err != nil {
			return
		}
	}
}

// handleError maps domain sentinel errors into standardized JSON error responses.
func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_AMOUNT",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_STATUS",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidPaymentMethod):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAYMENT_METHOD",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidDateFormat):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_DATE_FORMAT",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrPaymentNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "PAYMENT_NOT_FOUND",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrStudentNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "STUDENT_NOT_FOUND",
				"message": err.Error(),
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Error interno del servidor",
			},
		})
	}
}
