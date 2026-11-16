package schedule

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

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/schedule")
	handler.RegisterRoutes(group)
	return r
}

func TestHandlerAvailabilityRoutes(t *testing.T) {
	repo := newInMemoryRepo()
	reader := &mockStudentReader{students: make(map[string]StudentSummary)}
	svc := NewService(repo, reader)
	h := NewHandler(svc)
	r := setupTestRouter(h)

	// 1. GET /api/v1/schedule/availability (returns 7 days)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/schedule/availability", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var avails []DayAvailabilityResponse
	if err := json.Unmarshal(w.Body.Bytes(), &avails); err != nil {
		t.Fatalf("failed unmarshaling: %v", err)
	}
	if len(avails) != 7 {
		t.Errorf("expected 7 canonical days, got %d", len(avails))
	}

	// 2. PUT /api/v1/schedule/availability/Lun (valid interval)
	payload := map[string]interface{}{
		"enabled": true,
		"intervals": []map[string]string{
			{"startTime": "08:00", "endTime": "13:00"},
		},
	}
	body, _ := json.Marshal(payload)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/schedule/availability/Lun", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// 3. PUT /api/v1/schedule/availability/invalid_day -> 400 Bad Request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/schedule/availability/noday", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid day, got %d", w.Code)
	}
}

func TestHandlerSessionAndOccurrenceRoutes(t *testing.T) {
	repo := newInMemoryRepo()
	studentID := primitive.NewObjectID().Hex()
	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			studentID: {ID: studentID, FullName: "Alumno HTTP", Archived: false},
		},
	}
	svc := NewService(repo, reader)
	h := NewHandler(svc)
	r := setupTestRouter(h)

	// Configure Monday availability
	enabled := true
	_, _ = svc.UpdateAvailabilityForDay(context.Background(), "Lun", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "13:00"},
		},
	})

	// 1. POST /api/v1/schedule/sessions -> 201 Created
	payload := map[string]interface{}{
		"day":        "Lun",
		"startTime":  "08:00",
		"endTime":    "09:15",
		"startDate":  "2026-11-01",
		"studentIds": []string{studentID},
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/schedule/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created SessionResponse
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.ID == "" {
		t.Fatalf("expected valid session ID returned")
	}

	// 2. GET /api/v1/schedule/sessions/:id -> 200 OK
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/schedule/sessions/"+created.ID, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 getting session by ID, got %d", w.Code)
	}

	// 3. GET /api/v1/schedule/sessions?from=2026-11-02&to=2026-11-08 -> 200 OK with occurrences
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/schedule/sessions?from=2026-11-02&to=2026-11-08", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 getting occurrences, got %d", w.Code)
	}
	var occs []SessionOccurrenceResponse
	_ = json.Unmarshal(w.Body.Bytes(), &occs)
	if len(occs) != 1 {
		t.Errorf("expected 1 occurrence, got %d", len(occs))
	}

	// 4. DELETE /api/v1/schedule/sessions/:id -> 200 OK
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/schedule/sessions/"+created.ID, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 archiving session, got %d", w.Code)
	}
}

func TestHandlerEmptyScheduleReturnsArraysNotNil(t *testing.T) {
	repo := newInMemoryRepo()
	reader := &mockStudentReader{students: make(map[string]StudentSummary)}
	svc := NewService(repo, reader)
	h := NewHandler(svc)
	r := setupTestRouter(h)

	// 1. GET /api/v1/schedule/sessions with empty db -> Must be JSON '[]' (never 'null')
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/schedule/sessions", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	rawSessionsBody := strings.TrimSpace(w.Body.String())
	if rawSessionsBody == "null" {
		t.Fatalf("GET /api/v1/schedule/sessions returned literal 'null', expected '[]'")
	}
	if rawSessionsBody != "[]" {
		t.Fatalf("expected '[]', got '%s'", rawSessionsBody)
	}

	// 2. GET /api/v1/schedule/availability with unconfigured db -> intervals must be '[]' (never null)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/schedule/availability", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	rawAvailBody := w.Body.String()
	if strings.Contains(rawAvailBody, `"intervals":null`) {
		t.Fatalf("GET /api/v1/schedule/availability returned null intervals in: %s", rawAvailBody)
	}

	// 3. GET /api/v1/schedule/exceptions with empty db -> Must be JSON '[]' (never 'null')
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/schedule/exceptions", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	rawExceptionsBody := strings.TrimSpace(w.Body.String())
	if rawExceptionsBody == "null" {
		t.Fatalf("GET /api/v1/schedule/exceptions returned literal 'null', expected '[]'")
	}
	if rawExceptionsBody != "[]" {
		t.Fatalf("expected '[]', got '%s'", rawExceptionsBody)
	}
}
