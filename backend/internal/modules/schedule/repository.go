package schedule

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines all database operations for the Schedule module.
type Repository interface {
	// Availability
	GetAllAvailability(ctx context.Context) ([]DayAvailability, error)
	GetAvailabilityByDay(ctx context.Context, day DayOfWeek) (*DayAvailability, error)
	UpsertAvailability(ctx context.Context, availability *DayAvailability) error
	DeleteAvailabilityByDay(ctx context.Context, day DayOfWeek) error

	// Sessions
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id primitive.ObjectID) (*Session, error)
	ListActiveSessions(ctx context.Context, day *DayOfWeek, studentID *primitive.ObjectID) ([]Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	ArchiveSession(ctx context.Context, id primitive.ObjectID) error

	// Exceptions
	CreateException(ctx context.Context, exception *ScheduleException) error
	GetExceptionByID(ctx context.Context, id primitive.ObjectID) (*ScheduleException, error)
	GetExceptionBySessionAndDate(ctx context.Context, sessionID primitive.ObjectID, originalDate string) (*ScheduleException, error)
	ListExceptionsByPeriod(ctx context.Context, from, to string) ([]ScheduleException, error)
	ListExceptionsForSession(ctx context.Context, sessionID primitive.ObjectID) ([]ScheduleException, error)
	DeleteException(ctx context.Context, id primitive.ObjectID) error

	// Maintenance
	EnsureIndexes(ctx context.Context) error
}

type mongoRepository struct {
	availabilitiesColl *mongo.Collection
	sessionsColl       *mongo.Collection
	exceptionsColl     *mongo.Collection
}

// NewRepository creates a new MongoDB repository for the Schedule module.
func NewRepository(db *mongo.Database) Repository {
	if db == nil {
		return &mongoRepository{}
	}
	return &mongoRepository{
		availabilitiesColl: db.Collection("availabilities"),
		sessionsColl:       db.Collection("sessions"),
		exceptionsColl:     db.Collection("schedule_exceptions"),
	}
}

// EnsureIndexes creates necessary MongoDB indexes for performance and uniqueness.
func (r *mongoRepository) EnsureIndexes(ctx context.Context) error {
	if r.availabilitiesColl == nil || r.sessionsColl == nil || r.exceptionsColl == nil {
		return ErrDatabaseUnavailable
	}

	// 1. Availabilities: Unique index on day
	_, err := r.availabilitiesColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "day", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_day"),
	})
	if err != nil {
		return err
	}

	// 2. Sessions: Indexes on day + status, and student_ids + status
	_, err = r.sessionsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "day", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_day_status"),
		},
		{
			Keys:    bson.D{{Key: "student_ids", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_student_ids_status"),
		},
	})
	if err != nil {
		return err
	}

	// 3. Schedule Exceptions: Unique index on session_id + original_date
	_, err = r.exceptionsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "session_id", Value: 1}, {Key: "original_date", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_session_original_date"),
		},
		{
			Keys:    bson.D{{Key: "original_date", Value: 1}},
			Options: options.Index().SetName("idx_original_date"),
		},
		{
			Keys:    bson.D{{Key: "new_date", Value: 1}},
			Options: options.Index().SetName("idx_new_date"),
		},
	})
	return err
}

// ==========================================
// AVAILABILITY METHODS
// ==========================================

func (r *mongoRepository) GetAllAvailability(ctx context.Context) ([]DayAvailability, error) {
	if r.availabilitiesColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	cursor, err := r.availabilitiesColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []DayAvailability
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []DayAvailability{}
	}
	for i := range list {
		if list[i].Intervals == nil {
			list[i].Intervals = []AvailabilityInterval{}
		}
	}
	return list, nil
}

func (r *mongoRepository) GetAvailabilityByDay(ctx context.Context, day DayOfWeek) (*DayAvailability, error) {
	if r.availabilitiesColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	var avail DayAvailability
	err := r.availabilitiesColl.FindOne(ctx, bson.M{"day": day}).Decode(&avail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if avail.Intervals == nil {
		avail.Intervals = []AvailabilityInterval{}
	}
	return &avail, nil
}

func (r *mongoRepository) UpsertAvailability(ctx context.Context, availability *DayAvailability) error {
	if r.availabilitiesColl == nil {
		return ErrDatabaseUnavailable
	}

	now := time.Now().UTC()
	availability.UpdatedAt = now

	filter := bson.M{"day": availability.Day}
	update := bson.M{
		"$set": bson.M{
			"day":        availability.Day,
			"enabled":    availability.Enabled,
			"intervals":  availability.Intervals,
			"timezone":   availability.Timezone,
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	opts := options.Update().SetUpsert(true)
	res, err := r.availabilitiesColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.UpsertedID != nil {
		if oid, ok := res.UpsertedID.(primitive.ObjectID); ok {
			availability.ID = oid
		}
	}
	return nil
}

func (r *mongoRepository) DeleteAvailabilityByDay(ctx context.Context, day DayOfWeek) error {
	if r.availabilitiesColl == nil {
		return ErrDatabaseUnavailable
	}

	_, err := r.availabilitiesColl.DeleteOne(ctx, bson.M{"day": day})
	return err
}

// ==========================================
// SESSIONS METHODS
// ==========================================

func (r *mongoRepository) CreateSession(ctx context.Context, session *Session) error {
	if r.sessionsColl == nil {
		return ErrDatabaseUnavailable
	}

	now := time.Now().UTC()
	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	session.CreatedAt = now
	session.UpdatedAt = now
	if session.Status == "" {
		session.Status = SessionStatusActive
	}

	_, err := r.sessionsColl.InsertOne(ctx, session)
	return err
}

func (r *mongoRepository) GetSessionByID(ctx context.Context, id primitive.ObjectID) (*Session, error) {
	if r.sessionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	filter := bson.M{
		"_id":    id,
		"status": SessionStatusActive,
	}

	var session Session
	err := r.sessionsColl.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *mongoRepository) ListActiveSessions(ctx context.Context, day *DayOfWeek, studentID *primitive.ObjectID) ([]Session, error) {
	if r.sessionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	filter := bson.M{
		"status": SessionStatusActive,
	}
	if day != nil && *day != "" {
		filter["day"] = *day
	}
	if studentID != nil && !studentID.IsZero() {
		filter["student_ids"] = *studentID
	}

	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})
	cursor, err := r.sessionsColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []Session
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Session{}
	}
	for i := range list {
		if list[i].StudentIDs == nil {
			list[i].StudentIDs = []primitive.ObjectID{}
		}
	}
	return list, nil
}

func (r *mongoRepository) UpdateSession(ctx context.Context, session *Session) error {
	if r.sessionsColl == nil {
		return ErrDatabaseUnavailable
	}

	session.UpdatedAt = time.Now().UTC()
	filter := bson.M{
		"_id":    session.ID,
		"status": SessionStatusActive,
	}

	update := bson.M{
		"$set": bson.M{
			"day":              session.Day,
			"start_time":       session.StartTime,
			"end_time":         session.EndTime,
			"duration_minutes": session.DurationMinutes,
			"type":             session.Type,
			"student_ids":      session.StudentIDs,
			"is_recurring":     session.IsRecurring,
			"start_date":       session.StartDate,
			"end_date":         session.EndDate,
			"status":           session.Status,
			"notes":            session.Notes,
			"updated_at":       session.UpdatedAt,
		},
	}

	res, err := r.sessionsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (r *mongoRepository) ArchiveSession(ctx context.Context, id primitive.ObjectID) error {
	if r.sessionsColl == nil {
		return ErrDatabaseUnavailable
	}

	now := time.Now().UTC()
	filter := bson.M{
		"_id":    id,
		"status": SessionStatusActive,
	}
	update := bson.M{
		"$set": bson.M{
			"status":      SessionStatusArchived,
			"archived_at": now,
			"updated_at":  now,
		},
	}

	res, err := r.sessionsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrSessionNotFound
	}

	// Limpiar cualquier excepción asociada a la sesión archivada
	if r.exceptionsColl != nil {
		_, _ = r.exceptionsColl.DeleteMany(ctx, bson.M{"session_id": id})
	}

	return nil
}

// ==========================================
// SCHEDULE EXCEPTIONS METHODS
// ==========================================

func (r *mongoRepository) CreateException(ctx context.Context, exception *ScheduleException) error {
	if r.exceptionsColl == nil {
		return ErrDatabaseUnavailable
	}

	now := time.Now().UTC()
	if exception.ID.IsZero() {
		exception.ID = primitive.NewObjectID()
	}
	exception.CreatedAt = now
	exception.UpdatedAt = now

	_, err := r.exceptionsColl.InsertOne(ctx, exception)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateException
		}
		return err
	}
	return nil
}

func (r *mongoRepository) GetExceptionByID(ctx context.Context, id primitive.ObjectID) (*ScheduleException, error) {
	if r.exceptionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	var exc ScheduleException
	err := r.exceptionsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&exc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrExceptionNotFound
		}
		return nil, err
	}
	return &exc, nil
}

func (r *mongoRepository) GetExceptionBySessionAndDate(ctx context.Context, sessionID primitive.ObjectID, originalDate string) (*ScheduleException, error) {
	if r.exceptionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	var exc ScheduleException
	filter := bson.M{
		"session_id":    sessionID,
		"original_date": originalDate,
	}
	err := r.exceptionsColl.FindOne(ctx, filter).Decode(&exc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &exc, nil
}

func (r *mongoRepository) ListExceptionsByPeriod(ctx context.Context, from, to string) ([]ScheduleException, error) {
	if r.exceptionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	// An exception belongs to the period if either:
	// 1. original_date is in [from, to]
	// 2. new_date is in [from, to] (for rescheduled sessions moved into this period)
	filter := bson.M{
		"$or": []bson.M{
			{
				"original_date": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
			{
				"new_date": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
	}

	cursor, err := r.exceptionsColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []ScheduleException
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []ScheduleException{}
	}
	return list, nil
}

func (r *mongoRepository) ListExceptionsForSession(ctx context.Context, sessionID primitive.ObjectID) ([]ScheduleException, error) {
	if r.exceptionsColl == nil {
		return nil, ErrDatabaseUnavailable
	}

	cursor, err := r.exceptionsColl.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []ScheduleException
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []ScheduleException{}
	}
	return list, nil
}

func (r *mongoRepository) DeleteException(ctx context.Context, id primitive.ObjectID) error {
	if r.exceptionsColl == nil {
		return ErrDatabaseUnavailable
	}

	res, err := r.exceptionsColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrExceptionNotFound
	}
	return nil
}
