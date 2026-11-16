package students

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRouter() (*gin.Engine, Service) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := newMockRepository()
	svc := NewService(repo)
	handler := NewHandler(svc)

	v1 := r.Group("/api/v1/students")
	handler.RegisterRoutes(v1)

	return r, svc
}

func TestHandlerCreateSuccess(t *testing.T) {
	r, _ := setupTestRouter()

	body := []byte(`{"fullName":"Claudio Bravo","frequency":"2 veces/sem."}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("código esperado 201, obtenido %d. Cuerpo: %s", w.Code, w.Body.String())
	}

	var resp StudentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al decodificar respuesta JSON: %v", err)
	}
	if resp.FullName != "Claudio Bravo" {
		t.Errorf("nombre esperado 'Claudio Bravo', obtenido '%s'", resp.FullName)
	}
	if resp.Frequency != "2 veces/sem." {
		t.Errorf("frecuencia esperada '2 veces/sem.', obtenida '%s'", resp.Frequency)
	}
}

func TestHandlerCreateInvalidJSON(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBufferString("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("código esperado 400, obtenido %d", w.Code)
	}

	var errResp ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error.Code != "INVALID_JSON" {
		t.Errorf("código de error esperado 'INVALID_JSON', obtenido '%s'", errResp.Error.Code)
	}
}

func TestHandlerCreateEmptyName(t *testing.T) {
	r, _ := setupTestRouter()

	body := []byte(`{"fullName":"","frequency":"2 veces/sem."}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("código esperado 400, obtenido %d", w.Code)
	}

	var errResp ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error.Code != "EMPTY_NAME" {
		t.Errorf("código de error esperado 'EMPTY_NAME', obtenido '%s'", errResp.Error.Code)
	}
}

func TestHandlerGetByID(t *testing.T) {
	r, svc := setupTestRouter()

	created, _ := svc.Create(t.Context(), CreateStudentRequest{
		FullName:  "Alexis Sanchez",
		Frequency: Frequency3x,
	})

	// Petición válida
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+created.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	// ID inválido
	reqInv := httptest.NewRequest(http.MethodGet, "/api/v1/students/id-invalido", nil)
	wInv := httptest.NewRecorder()
	r.ServeHTTP(wInv, reqInv)
	if wInv.Code != http.StatusBadRequest {
		t.Fatalf("código esperado 400 para ID inválido, obtenido %d", wInv.Code)
	}

	// ID inexistente
	nonExistentHex := primitive.NewObjectID().Hex()
	req404 := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+nonExistentHex, nil)
	w404 := httptest.NewRecorder()
	r.ServeHTTP(w404, req404)
	if w404.Code != http.StatusNotFound {
		t.Fatalf("código esperado 404 para ID inexistente, obtenido %d", w404.Code)
	}
}

func TestHandlerList(t *testing.T) {
	r, svc := setupTestRouter()

	_, _ = svc.Create(t.Context(), CreateStudentRequest{
		FullName:  "Alumno 1",
		Frequency: Frequency1x,
	})
	_, _ = svc.Create(t.Context(), CreateStudentRequest{
		FullName:  "Alumno 2",
		Frequency: Frequency2x,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	var list []StudentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("error al deserializar lista: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("longitud esperada 2, obtenida %d", len(list))
	}
}

func TestHandlerUpdate(t *testing.T) {
	r, svc := setupTestRouter()

	created, _ := svc.Create(t.Context(), CreateStudentRequest{
		FullName:  "Nombre Inicial",
		Frequency: Frequency1x,
	})

	body := []byte(`{"fullName":"Nombre Editado","frequency":"3 veces/sem."}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/students/"+created.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	var updated StudentResponse
	_ = json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.FullName != "Nombre Editado" {
		t.Errorf("nombre esperado 'Nombre Editado', obtenido '%s'", updated.FullName)
	}
}

func TestHandlerDelete(t *testing.T) {
	r, svc := setupTestRouter()

	created, _ := svc.Create(t.Context(), CreateStudentRequest{
		FullName:  "Alumno Para Eliminar",
		Frequency: Frequency2x,
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/"+created.ID, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("código esperado 200, obtenido %d", w.Code)
	}

	// Comprobación de que no se encuentra tras borrado
	reqCheck := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+created.ID, nil)
	wCheck := httptest.NewRecorder()
	r.ServeHTTP(wCheck, reqCheck)
	if wCheck.Code != http.StatusNotFound {
		t.Errorf("código esperado 404 tras eliminación, obtenido %d", wCheck.Code)
	}
}
