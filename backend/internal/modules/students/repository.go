package students

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionName = "students"

// Repository define el contrato de persistencia exclusivo para la entidad Student.
type Repository interface {
	Create(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Student, error)
	List(ctx context.Context, filter StudentListFilter) ([]Student, error)
	Update(ctx context.Context, student *Student) error
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
}

type mongoRepository struct {
	collection *mongo.Collection
}

// NewRepository inicializa una nueva instancia del repositorio MongoDB para alumnos.
func NewRepository(db *mongo.Database) Repository {
	var coll *mongo.Collection
	if db != nil {
		coll = db.Collection(CollectionName)
	}
	return &mongoRepository{
		collection: coll,
	}
}

// Create inserta un nuevo alumno en la colección 'students'.
func (r *mongoRepository) Create(ctx context.Context, student *Student) error {
	if r.collection == nil {
		return ErrDatabaseUnavailable
	}

	if student.ID.IsZero() {
		student.ID = primitive.NewObjectID()
	}

	now := time.Now().UTC()
	student.CreatedAt = now
	student.UpdatedAt = now
	student.ArchivedAt = nil

	_, err := r.collection.InsertOne(ctx, student)
	if err != nil {
		return fmt.Errorf("error al insertar alumno en MongoDB: %w", err)
	}

	return nil
}

// GetByID consulta un alumno por su ObjectID, asegurando que no se encuentre archivado.
func (r *mongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Student, error) {
	if r.collection == nil {
		return nil, ErrDatabaseUnavailable
	}

	filter := bson.M{
		"_id":        id,
		"archivedAt": nil,
	}

	var student Student
	err := r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("error al buscar alumno por ID: %w", err)
	}

	return &student, nil
}

// List recupera todos los alumnos no archivados, admitiendo filtros opcionales de nombre, estado y frecuencia.
func (r *mongoRepository) List(ctx context.Context, filter StudentListFilter) ([]Student, error) {
	if r.collection == nil {
		return nil, ErrDatabaseUnavailable
	}

	filterDoc := bson.M{
		"archivedAt": nil,
	}

	if filter.Name != "" {
		// Búsqueda case-insensitive por subcadena escapada
		escaped := regexp.QuoteMeta(filter.Name)
		filterDoc["fullName"] = primitive.Regex{Pattern: escaped, Options: "i"}
	}

	if filter.Status != "" {
		filterDoc["membershipStatus"] = filter.Status
	}

	if filter.Frequency != "" {
		filterDoc["frequency"] = filter.Frequency
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filterDoc, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error al listar alumnos en MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	students := make([]Student, 0)
	if err := cursor.All(ctx, &students); err != nil {
		return nil, fmt.Errorf("error al decodificar lista de alumnos: %w", err)
	}

	return students, nil
}

// Update actualiza los datos de un alumno existente y renueva su timestamp 'updatedAt'.
func (r *mongoRepository) Update(ctx context.Context, student *Student) error {
	if r.collection == nil {
		return ErrDatabaseUnavailable
	}

	student.UpdatedAt = time.Now().UTC()

	filter := bson.M{
		"_id":        student.ID,
		"archivedAt": nil,
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"fullName":            student.FullName,
			"email":               student.Email,
			"phone":               student.Phone,
			"planName":            student.PlanName,
			"frequency":           student.Frequency,
			"membershipStatus":    student.MembershipStatus,
			"fixedSchedule":       student.FixedSchedule,
			"membershipStartDate": student.MembershipStartDate,
			"membershipDueDate":   student.MembershipDueDate,
			"notes":               student.Notes,
			"healthConditions":    student.HealthConditions,
			"emergencyContact":    student.EmergencyContact,
			"updatedAt":           student.UpdatedAt,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("error al actualizar alumno en MongoDB: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrStudentNotFound
	}

	return nil
}

// SoftDelete marca lógicamente un alumno con archivedAt y preserva sus datos históricos.
func (r *mongoRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	if r.collection == nil {
		return ErrDatabaseUnavailable
	}

	now := time.Now().UTC()
	filter := bson.M{
		"_id":        id,
		"archivedAt": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"archivedAt": now,
			"updatedAt":  now,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error al archivar alumno en MongoDB: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrStudentNotFound
	}

	return nil
}
