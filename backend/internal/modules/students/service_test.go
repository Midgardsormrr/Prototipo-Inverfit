package students

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	students map[primitive.ObjectID]*Student
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		students: make(map[primitive.ObjectID]*Student),
	}
}

func (m *mockRepository) Create(ctx context.Context, student *Student) error {
	if student.ID.IsZero() {
		student.ID = primitive.NewObjectID()
	}
	now := time.Now().UTC()
	student.CreatedAt = now
	student.UpdatedAt = now
	m.students[student.ID] = student
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Student, error) {
	s, ok := m.students[id]
	if !ok || s.ArchivedAt != nil {
		return nil, ErrStudentNotFound
	}
	return s, nil
}

func (m *mockRepository) List(ctx context.Context, filter StudentListFilter) ([]Student, error) {
	result := make([]Student, 0)
	for _, s := range m.students {
		if s.ArchivedAt != nil {
			continue
		}
		if filter.Status != "" && s.MembershipStatus != filter.Status {
			continue
		}
		if filter.Frequency != "" && s.Frequency != filter.Frequency {
			continue
		}
		result = append(result, *s)
	}
	return result, nil
}

func (m *mockRepository) Update(ctx context.Context, student *Student) error {
	s, ok := m.students[student.ID]
	if !ok || s.ArchivedAt != nil {
		return ErrStudentNotFound
	}
	student.UpdatedAt = time.Now().UTC()
	m.students[student.ID] = student
	return nil
}

func (m *mockRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	s, ok := m.students[id]
	if !ok || s.ArchivedAt != nil {
		return ErrStudentNotFound
	}
	now := time.Now().UTC()
	s.ArchivedAt = &now
	s.UpdatedAt = now
	return nil
}

func TestServiceCreateValid(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	req := CreateStudentRequest{
		FullName:  "Carlos Valenzuela",
		Frequency: Frequency2x,
	}

	res, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("error inesperado al crear alumno: %v", err)
	}

	if res.FullName != "Carlos Valenzuela" {
		t.Errorf("nombre esperado 'Carlos Valenzuela', obtenido '%s'", res.FullName)
	}
	if res.Frequency != Frequency2x {
		t.Errorf("frecuencia esperada '%s', obtenida '%s'", Frequency2x, res.Frequency)
	}
	if res.MembershipStatus != StatusActive {
		t.Errorf("estado por defecto esperado '%s', obtenido '%s'", StatusActive, res.MembershipStatus)
	}
	if res.PlanName != "Membresía Mensual" {
		t.Errorf("plan por defecto esperado 'Membresía Mensual', obtenido '%s'", res.PlanName)
	}
	if res.ID == "" {
		t.Errorf("ID generado no debe estar vacío")
	}
}

func TestServiceCreateEmptyName(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	req := CreateStudentRequest{
		FullName:  "   ",
		Frequency: Frequency2x,
	}

	_, err := svc.Create(context.Background(), req)
	if err != ErrEmptyName {
		t.Fatalf("error esperado ErrEmptyName, obtenido: %v", err)
	}
}

func TestServiceCreateInvalidFrequency(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	req := CreateStudentRequest{
		FullName:  "Valeria Gomez",
		Frequency: "5 veces/sem.",
	}

	_, err := svc.Create(context.Background(), req)
	if err != ErrInvalidFrequency {
		t.Fatalf("error esperado ErrInvalidFrequency, obtenido: %v", err)
	}
}

func TestServiceCreateInvalidStatus(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	req := CreateStudentRequest{
		FullName:         "Valeria Gomez",
		Frequency:        Frequency2x,
		MembershipStatus: "Suspendido",
	}

	_, err := svc.Create(context.Background(), req)
	if err != ErrInvalidStatus {
		t.Fatalf("error esperado ErrInvalidStatus, obtenido: %v", err)
	}
}

func TestServiceGetByID(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	created, err := svc.Create(context.Background(), CreateStudentRequest{
		FullName:  "Andres Bello",
		Frequency: Frequency3x,
	})
	if err != nil {
		t.Fatalf("error al crear: %v", err)
	}

	// Consulta válida
	found, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("error al consultar por ID: %v", err)
	}
	if found.FullName != "Andres Bello" {
		t.Errorf("nombre esperado 'Andres Bello', obtenido '%s'", found.FullName)
	}

	// ID inválido
	_, err = svc.GetByID(context.Background(), "no-es-un-hex")
	if err != ErrInvalidID {
		t.Errorf("error esperado ErrInvalidID, obtenido: %v", err)
	}

	// ID inexistente
	nonExistentHex := primitive.NewObjectID().Hex()
	_, err = svc.GetByID(context.Background(), nonExistentHex)
	if err != ErrStudentNotFound {
		t.Errorf("error esperado ErrStudentNotFound, obtenido: %v", err)
	}
}

func TestServiceListAndFilter(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	_, _ = svc.Create(context.Background(), CreateStudentRequest{
		FullName:         "Alumno Activo 1",
		Frequency:        Frequency2x,
		MembershipStatus: StatusActive,
	})
	_, _ = svc.Create(context.Background(), CreateStudentRequest{
		FullName:         "Alumno Inactivo 1",
		Frequency:        Frequency1x,
		MembershipStatus: StatusInactive,
	})

	// Listar todos
	all, err := svc.List(context.Background(), StudentListFilter{})
	if err != nil {
		t.Fatalf("error al listar: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("total esperado 2, obtenido %d", len(all))
	}

	// Filtrar por estado
	actives, err := svc.List(context.Background(), StudentListFilter{Status: StatusActive})
	if err != nil {
		t.Fatalf("error al filtrar por estado: %v", err)
	}
	if len(actives) != 1 || actives[0].FullName != "Alumno Activo 1" {
		t.Errorf("filtro de activos falló")
	}

	// Filtrar con estado inválido
	_, err = svc.List(context.Background(), StudentListFilter{Status: "EstadoInvalido"})
	if err != ErrInvalidStatus {
		t.Errorf("error esperado ErrInvalidStatus, obtenido: %v", err)
	}
}

func TestServiceUpdate(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	created, _ := svc.Create(context.Background(), CreateStudentRequest{
		FullName:  "Nombre Original",
		Frequency: Frequency2x,
	})

	newName := "Nombre Modificado"
	newFreq := Frequency3x
	updated, err := svc.Update(context.Background(), created.ID, UpdateStudentRequest{
		FullName:  &newName,
		Frequency: &newFreq,
	})
	if err != nil {
		t.Fatalf("error al actualizar: %v", err)
	}
	if updated.FullName != "Nombre Modificado" {
		t.Errorf("nombre esperado 'Nombre Modificado', obtenido '%s'", updated.FullName)
	}
	if updated.Frequency != Frequency3x {
		t.Errorf("frecuencia esperada '%s', obtenida '%s'", Frequency3x, updated.Frequency)
	}

	// Modificación con nombre vacío
	emptyName := "   "
	_, err = svc.Update(context.Background(), created.ID, UpdateStudentRequest{
		FullName: &emptyName,
	})
	if err != ErrEmptyName {
		t.Errorf("error esperado ErrEmptyName, obtenido: %v", err)
	}
}

func TestServiceDelete(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	created, _ := svc.Create(context.Background(), CreateStudentRequest{
		FullName:  "Alumno Para Archivar",
		Frequency: Frequency2x,
	})

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("error al archivar alumno: %v", err)
	}

	// Ya no debe ser accesible por GetByID
	_, err := svc.GetByID(context.Background(), created.ID)
	if err != ErrStudentNotFound {
		t.Errorf("alumno archivado no debe ser consultable, error obtenido: %v", err)
	}

	// Tampoco debe figurar en List
	list, _ := svc.List(context.Background(), StudentListFilter{})
	if len(list) != 0 {
		t.Errorf("alumno archivado no debe aparecer en listado activo, obtenido %d", len(list))
	}
}
