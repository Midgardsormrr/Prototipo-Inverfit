package schedule

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"inverfit-backend/internal/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getScheduleTestMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	candidates := []string{"../../../.env", "../../.env", ".env", "backend/.env"}
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

func TestRepositoryNilDatabase(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	_, err := repo.GetAllAvailability(ctx)
	if err != ErrDatabaseUnavailable {
		t.Errorf("expected ErrDatabaseUnavailable, got: %v", err)
	}

	err = repo.CreateSession(ctx, &Session{})
	if err != ErrDatabaseUnavailable {
		t.Errorf("expected ErrDatabaseUnavailable, got: %v", err)
	}

	err = repo.CreateException(ctx, &ScheduleException{})
	if err != ErrDatabaseUnavailable {
		t.Errorf("expected ErrDatabaseUnavailable, got: %v", err)
	}
}

func TestRepositoryIntegrationWithMongoDB(t *testing.T) {
	uri := getScheduleTestMongoURI()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, uri, "inverfit_test")
	if err != nil {
		t.Skipf("MongoDB no disponible en %s para prueba de integración: %v", uri, err)
		return
	}
	defer func() {
		_ = client.Database.Collection("availabilities").Drop(context.Background())
		_ = client.Database.Collection("sessions").Drop(context.Background())
		_ = client.Database.Collection("schedule_exceptions").Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}()

	repo := NewRepository(client.Database)
	testCtx := context.Background()

	// 1. Ensure indexes
	if err := repo.EnsureIndexes(testCtx); err != nil {
		t.Fatalf("failed ensuring indexes: %v", err)
	}

	// 2. Upsert and Get Availability
	avail := &DayAvailability{
		Day:     DayLun,
		Enabled: true,
		Intervals: []AvailabilityInterval{
			{ID: "int-lun-1", StartTime: "08:00", EndTime: "12:00"},
		},
		Timezone: DefaultTimezone,
	}
	if err := repo.UpsertAvailability(testCtx, avail); err != nil {
		t.Fatalf("failed upserting availability: %v", err)
	}

	fetchedAvail, err := repo.GetAvailabilityByDay(testCtx, DayLun)
	if err != nil || fetchedAvail == nil {
		t.Fatalf("failed getting availability by day: %v", err)
	}
	if len(fetchedAvail.Intervals) != 1 || fetchedAvail.Intervals[0].StartTime != "08:00" {
		t.Errorf("unexpected intervals: %+v", fetchedAvail.Intervals)
	}

	// 3. Create, Get, Update, Archive Session
	studentID := primitive.NewObjectID()
	sess := &Session{
		Day:             DayLun,
		StartTime:       "08:30",
		EndTime:         "09:45",
		DurationMinutes: 75,
		Type:            SessionTypeIndividual,
		StudentIDs:      []primitive.ObjectID{studentID},
		IsRecurring:     true,
		StartDate:       "2026-11-02",
		Status:          SessionStatusActive,
		Notes:           "Prueba integración",
	}
	if err := repo.CreateSession(testCtx, sess); err != nil {
		t.Fatalf("failed creating session in mongo: %v", err)
	}
	if sess.ID.IsZero() {
		t.Fatalf("expected non-zero ID assigned to session")
	}

	fetchedSess, err := repo.GetSessionByID(testCtx, sess.ID)
	if err != nil {
		t.Fatalf("failed fetching session by ID: %v", err)
	}
	if fetchedSess.StartTime != "08:30" {
		t.Errorf("expected startTime 08:30, got %s", fetchedSess.StartTime)
	}

	// 4. Create and List Exception
	exc := &ScheduleException{
		SessionID:    sess.ID,
		OriginalDate: "2026-11-02",
		Type:         ExceptionTypeCancelled,
		Notes:        "Cancelación test",
	}
	if err := repo.CreateException(testCtx, exc); err != nil {
		t.Fatalf("failed creating exception: %v", err)
	}

	exceptions, err := repo.ListExceptionsByPeriod(testCtx, "2026-11-01", "2026-11-30")
	if err != nil {
		t.Fatalf("failed listing exceptions: %v", err)
	}
	if len(exceptions) != 1 {
		t.Errorf("expected 1 exception, got %d", len(exceptions))
	}

	// 5. Delete Exception
	if err := repo.DeleteException(testCtx, exc.ID); err != nil {
		t.Fatalf("failed deleting exception: %v", err)
	}

	// 6. Archive Session
	if err := repo.ArchiveSession(testCtx, sess.ID); err != nil {
		t.Fatalf("failed archiving session: %v", err)
	}

	_, err = repo.GetSessionByID(testCtx, sess.ID)
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("expected ErrSessionNotFound for archived session, got: %v", err)
	}
}
