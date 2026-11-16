package students

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"inverfit-backend/internal/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepositoryNilDatabase(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	dummyStudent := &Student{FullName: "Test"}

	if err := repo.Create(ctx, dummyStudent); err != ErrDatabaseUnavailable {
		t.Errorf("esperado ErrDatabaseUnavailable en Create, obtenido %v", err)
	}

	if _, err := repo.GetByID(ctx, primitive.NewObjectID()); err != ErrDatabaseUnavailable {
		t.Errorf("esperado ErrDatabaseUnavailable en GetByID, obtenido %v", err)
	}

	if _, err := repo.List(ctx, StudentListFilter{}); err != ErrDatabaseUnavailable {
		t.Errorf("esperado ErrDatabaseUnavailable en List, obtenido %v", err)
	}

	if err := repo.Update(ctx, dummyStudent); err != ErrDatabaseUnavailable {
		t.Errorf("esperado ErrDatabaseUnavailable en Update, obtenido %v", err)
	}

	if err := repo.SoftDelete(ctx, primitive.NewObjectID()); err != ErrDatabaseUnavailable {
		t.Errorf("esperado ErrDatabaseUnavailable en SoftDelete, obtenido %v", err)
	}
}

func getTestMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	candidates := []string{"../../.env", "../../../.env", ".env", "backend/.env"}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "MONGODB_URI=") {
				return strings.TrimSpace(strings.TrimPrefix(line, "MONGODB_URI="))
			}
		}
	}
	return ""
}

// TestRepositoryIntegrationWithMongoDB se conecta a la base de datos aislada 'inverfit_test'.
// Si MongoDB no está accesible (por ejemplo en entornos CI sin base de datos), el test se omite.
func TestRepositoryIntegrationWithMongoDB(t *testing.T) {
	mongoURI := getTestMongoURI()
	if mongoURI == "" {
		t.Skip("MONGODB_URI no configurado, omitiendo prueba de integración con MongoDB real")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testDBName := "inverfit_test"
	dbConn, err := database.Connect(ctx, mongoURI, testDBName)
	if err != nil {
		t.Skipf("No se pudo conectar a MongoDB en %s: %v. Omitiendo prueba de integración.", mongoURI, err)
	}
	defer dbConn.Disconnect(context.Background())

	coll := dbConn.Database.Collection(CollectionName)
	// Limpieza previa exclusiva de la colección en 'inverfit_test'
	_ = coll.Drop(ctx)
	defer coll.Drop(context.Background())

	repo := NewRepository(dbConn.Database)

	// 1. Create
	newStudent := &Student{
		FullName:         "Alumno Integracion",
		Frequency:        Frequency3x,
		MembershipStatus: StatusActive,
		PlanName:         "Plan personalizado",
		Email:            "integracion@inverfit.local",
	}

	if err := repo.Create(ctx, newStudent); err != nil {
		t.Fatalf("error en Create de repositorio: %v", err)
	}
	if newStudent.ID.IsZero() {
		t.Fatalf("ObjectID no fue asignado")
	}

	// 2. GetByID
	found, err := repo.GetByID(ctx, newStudent.ID)
	if err != nil {
		t.Fatalf("error en GetByID de repositorio: %v", err)
	}
	if found.FullName != "Alumno Integracion" {
		t.Errorf("nombre esperado 'Alumno Integracion', obtenido '%s'", found.FullName)
	}

	// 3. List
	list, err := repo.List(ctx, StudentListFilter{})
	if err != nil {
		t.Fatalf("error en List de repositorio: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("se esperaba 1 alumno en listado, obtenidos %d", len(list))
	}

	// 4. Update
	found.FullName = "Alumno Integracion Actualizado"
	found.Frequency = Frequency4x
	if err := repo.Update(ctx, found); err != nil {
		t.Fatalf("error en Update de repositorio: %v", err)
	}

	updated, err := repo.GetByID(ctx, newStudent.ID)
	if err != nil {
		t.Fatalf("error al re-consultar tras update: %v", err)
	}
	if updated.FullName != "Alumno Integracion Actualizado" {
		t.Errorf("nombre no actualizado, obtenido '%s'", updated.FullName)
	}

	// 5. SoftDelete
	if err := repo.SoftDelete(ctx, newStudent.ID); err != nil {
		t.Fatalf("error en SoftDelete de repositorio: %v", err)
	}

	// 6. Verificar que ya no aparece en GetByID activo
	_, err = repo.GetByID(ctx, newStudent.ID)
	if err != ErrStudentNotFound {
		t.Errorf("se esperaba ErrStudentNotFound tras SoftDelete, obtenido: %v", err)
	}

	// 7. Verificar que no aparece en List activo
	listAfterDelete, err := repo.List(ctx, StudentListFilter{})
	if err != nil {
		t.Fatalf("error en List tras SoftDelete: %v", err)
	}
	if len(listAfterDelete) != 0 {
		t.Errorf("el alumno archivado no debe figurar en List, encontrados: %d", len(listAfterDelete))
	}

	// 8. Verificar que el documento físico SIGUE existiendo en la colección con archivedAt poblado
	var rawDoc Student
	err = coll.FindOne(ctx, primitive.M{"_id": newStudent.ID}).Decode(&rawDoc)
	if err != nil {
		t.Fatalf("el documento físico fue eliminado erróneamente de la colección: %v", err)
	}
	if rawDoc.ArchivedAt == nil {
		t.Errorf("el campo archivedAt debe estar registrado")
	}
}
