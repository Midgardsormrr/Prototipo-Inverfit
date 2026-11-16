package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRouter() (*gin.Engine, Service, *mockStudentReader) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := newMockRepo()
	reader := newMockStudentReader()
	svc := NewService(repo, reader)
	h := NewHandler(svc)

	rg := r.Group("/api/v1/payments")
	h.RegisterRoutes(rg)

	return r, svc, reader
}

func TestHandlerEndpoints(t *testing.T) {
	router, svc, reader := setupTestRouter()

	// 1. Create a student in mock reader and a payment via service
	stOID := primitive.NewObjectID()
	studentID := stOID.Hex()
	reader.students[studentID] = &StudentInfo{
		ID:               studentID,
		FullName:         "Alumno Handler Test",
		PlanName:         "Plan Estándar",
		Frequency:        "2 veces/sem.",
		MembershipStatus: "Activo",
	}

	amt := int64(45000)
	resp, err := svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID:     studentID,
		Amount:        &amt,
		PaymentMethod: MethodTransfer,
		PaymentDate:   "2026-11-06",
		Notes:         "Test handler",
	})
	if err != nil {
		t.Fatalf("setup payment failed: %v", err)
	}

	// 2. GET /api/v1/payments
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/payments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /payments returned %d, want 200", w.Code)
	}

	// 3. GET /api/v1/payments/summary
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/payments/summary", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /payments/summary returned %d, want 200", w.Code)
	}

	// 4. GET /api/v1/payments/export
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/payments/export", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /payments/export returned %d, want 200", w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/csv") {
		t.Errorf("expected text/csv Content-Type, got %s", contentType)
	}
	bodyStr := w.Body.String()
	if !strings.Contains(bodyStr, "ID,Alumno,Plan") {
		t.Errorf("expected CSV header row, got %s", bodyStr)
	}

	// 5. GET /api/v1/payments/:id
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/payments/"+resp.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /payments/:id returned %d, want 200", w.Code)
	}

	// 6. PUT /api/v1/payments/:id
	newAmt := int64(48000)
	putBody, _ := json.Marshal(UpdatePaymentRequest{Amount: &newAmt})
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/payments/"+resp.ID, bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("PUT /payments/:id returned %d, want 200", w.Code)
	}

	// 7. POST /api/v1/payments/:id/simulate-reminder
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/payments/"+resp.ID+"/simulate-reminder", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("POST /simulate-reminder returned %d, want 200", w.Code)
	}

	// 8. POST /api/v1/payments/:id/simulate-receipt
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/payments/"+resp.ID+"/simulate-receipt", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("POST /simulate-receipt returned %d, want 200", w.Code)
	}

	// 9. DELETE /api/v1/payments/:id
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/payments/"+resp.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("DELETE /payments/:id returned %d, want 200", w.Code)
	}

	// 10. Verify 404 after deletion
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/payments/"+resp.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("GET /payments/:id after delete returned %d, want 404", w.Code)
	}
}
