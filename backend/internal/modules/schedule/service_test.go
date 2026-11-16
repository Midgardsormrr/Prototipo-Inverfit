package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// inMemoryRepository implements Repository for unit testing without MongoDB.
type inMemoryRepository struct {
	availabilities map[DayOfWeek]DayAvailability
	sessions       map[primitive.ObjectID]Session
	exceptions     map[primitive.ObjectID]ScheduleException
}

func newInMemoryRepo() *inMemoryRepository {
	return &inMemoryRepository{
		availabilities: make(map[DayOfWeek]DayAvailability),
		sessions:       make(map[primitive.ObjectID]Session),
		exceptions:     make(map[primitive.ObjectID]ScheduleException),
	}
}

func (m *inMemoryRepository) GetAllAvailability(ctx context.Context) ([]DayAvailability, error) {
	var list []DayAvailability
	for _, a := range m.availabilities {
		list = append(list, a)
	}
	return list, nil
}

func (m *inMemoryRepository) GetAvailabilityByDay(ctx context.Context, day DayOfWeek) (*DayAvailability, error) {
	if a, ok := m.availabilities[day]; ok {
		return &a, nil
	}
	return nil, nil
}

func (m *inMemoryRepository) UpsertAvailability(ctx context.Context, availability *DayAvailability) error {
	m.availabilities[availability.Day] = *availability
	return nil
}

func (m *inMemoryRepository) DeleteAvailabilityByDay(ctx context.Context, day DayOfWeek) error {
	delete(m.availabilities, day)
	return nil
}

func (m *inMemoryRepository) CreateSession(ctx context.Context, session *Session) error {
	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	m.sessions[session.ID] = *session
	return nil
}

func (m *inMemoryRepository) GetSessionByID(ctx context.Context, id primitive.ObjectID) (*Session, error) {
	if s, ok := m.sessions[id]; ok && s.Status == SessionStatusActive {
		return &s, nil
	}
	return nil, ErrSessionNotFound
}

func (m *inMemoryRepository) ListActiveSessions(ctx context.Context, day *DayOfWeek, studentID *primitive.ObjectID) ([]Session, error) {
	var list []Session
	for _, s := range m.sessions {
		if s.Status != SessionStatusActive {
			continue
		}
		if day != nil && s.Day != *day {
			continue
		}
		if studentID != nil {
			found := false
			for _, sid := range s.StudentIDs {
				if sid == *studentID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		list = append(list, s)
	}
	return list, nil
}

func (m *inMemoryRepository) UpdateSession(ctx context.Context, session *Session) error {
	if _, ok := m.sessions[session.ID]; !ok {
		return ErrSessionNotFound
	}
	m.sessions[session.ID] = *session
	return nil
}

func (m *inMemoryRepository) ArchiveSession(ctx context.Context, id primitive.ObjectID) error {
	s, ok := m.sessions[id]
	if !ok || s.Status != SessionStatusActive {
		return ErrSessionNotFound
	}
	s.Status = SessionStatusArchived
	now := time.Now().UTC()
	s.ArchivedAt = &now
	m.sessions[id] = s
	for excID, exc := range m.exceptions {
		if exc.SessionID == id {
			delete(m.exceptions, excID)
		}
	}
	return nil
}

func (m *inMemoryRepository) CreateException(ctx context.Context, exception *ScheduleException) error {
	for _, e := range m.exceptions {
		if e.SessionID == exception.SessionID && e.OriginalDate == exception.OriginalDate {
			return ErrDuplicateException
		}
	}
	if exception.ID.IsZero() {
		exception.ID = primitive.NewObjectID()
	}
	m.exceptions[exception.ID] = *exception
	return nil
}

func (m *inMemoryRepository) GetExceptionByID(ctx context.Context, id primitive.ObjectID) (*ScheduleException, error) {
	if e, ok := m.exceptions[id]; ok {
		return &e, nil
	}
	return nil, ErrExceptionNotFound
}

func (m *inMemoryRepository) GetExceptionBySessionAndDate(ctx context.Context, sessionID primitive.ObjectID, originalDate string) (*ScheduleException, error) {
	for _, e := range m.exceptions {
		if e.SessionID == sessionID && e.OriginalDate == originalDate {
			return &e, nil
		}
	}
	return nil, nil
}

func (m *inMemoryRepository) ListExceptionsByPeriod(ctx context.Context, from, to string) ([]ScheduleException, error) {
	var list []ScheduleException
	for _, e := range m.exceptions {
		if (e.OriginalDate >= from && e.OriginalDate <= to) || (e.NewDate != nil && *e.NewDate >= from && *e.NewDate <= to) {
			list = append(list, e)
		}
	}
	return list, nil
}

func (m *inMemoryRepository) ListExceptionsForSession(ctx context.Context, sessionID primitive.ObjectID) ([]ScheduleException, error) {
	var list []ScheduleException
	for _, e := range m.exceptions {
		if e.SessionID == sessionID {
			list = append(list, e)
		}
	}
	return list, nil
}

func (m *inMemoryRepository) DeleteException(ctx context.Context, id primitive.ObjectID) error {
	if _, ok := m.exceptions[id]; !ok {
		return ErrExceptionNotFound
	}
	delete(m.exceptions, id)
	return nil
}

func (m *inMemoryRepository) EnsureIndexes(ctx context.Context) error {
	return nil
}

// mockStudentReader simulates StudentReader for unit tests.
type mockStudentReader struct {
	students map[string]StudentSummary
}

func (m *mockStudentReader) GetStudentsSummary(ctx context.Context, ids []string) ([]StudentSummary, error) {
	var res []StudentSummary
	for _, id := range ids {
		if s, ok := m.students[id]; ok {
			res = append(res, s)
		}
	}
	return res, nil
}

// ==========================================
// TESTS
// ==========================================

func TestAvailabilityManagement(t *testing.T) {
	repo := newInMemoryRepo()
	reader := &mockStudentReader{students: make(map[string]StudentSummary)}
	svc := NewService(repo, reader)
	ctx := context.Background()

	// 1. Configure availability for Tuesday (Mar) with 2 intervals
	enabled := true
	req := UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "12:00"},
			{StartTime: "15:00", EndTime: "19:00"},
		},
	}

	avail, err := svc.UpdateAvailabilityForDay(ctx, "Mar", req)
	if err != nil {
		t.Fatalf("expected no error updating availability, got: %v", err)
	}
	if avail.Day != DayMar || len(avail.Intervals) != 2 {
		t.Errorf("unexpected availability returned: %+v", avail)
	}

	// 2. Reject overlapping intervals on same day
	overlapReq := UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "12:00"},
			{StartTime: "11:00", EndTime: "15:00"}, // overlaps 11:00 - 12:00
		},
	}
	_, err = svc.UpdateAvailabilityForDay(ctx, "Mar", overlapReq)
	if !errors.Is(err, ErrOverlappingIntervals) {
		t.Errorf("expected ErrOverlappingIntervals, got: %v", err)
	}

	// 3. Allow consecutive intervals on same day
	consecReq := UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "12:00"},
			{StartTime: "12:00", EndTime: "16:00"}, // consecutive: valid!
		},
	}
	avail, err = svc.UpdateAvailabilityForDay(ctx, "Mar", consecReq)
	if err != nil {
		t.Errorf("expected consecutive intervals to be accepted, got error: %v", err)
	}

	// 4. Reject outside operating hours (06:30 - 22:00)
	outsideReq := UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "05:00", EndTime: "08:00"},
		},
	}
	_, err = svc.UpdateAvailabilityForDay(ctx, "Mar", outsideReq)
	if !errors.Is(err, ErrOutsideOperatingHours) {
		t.Errorf("expected ErrOutsideOperatingHours, got: %v", err)
	}
}

func TestSessionCreationAndConsecutiveSessions(t *testing.T) {
	repo := newInMemoryRepo()
	student1ID := primitive.NewObjectID().Hex()
	student2ID := primitive.NewObjectID().Hex()
	archivedID := primitive.NewObjectID().Hex()

	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			student1ID: {ID: student1ID, FullName: "Estudiante Uno", Archived: false},
			student2ID: {ID: student2ID, FullName: "Estudiante Dos", Archived: false},
			archivedID: {ID: archivedID, FullName: "Estudiante Archivada", Archived: true},
		},
	}
	svc := NewService(repo, reader)
	ctx := context.Background()

	// Setup coach availability for Wednesday (Mié) 08:00 - 13:00
	enabled := true
	_, err := svc.UpdateAvailabilityForDay(ctx, "Mié", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "13:00"},
		},
	})
	if err != nil {
		t.Fatalf("failed setting availability: %v", err)
	}

	// 1. Create first session: 08:00 - 09:15 (75 min)
	s1, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "08:00",
		EndTime:    "09:15",
		StudentIDs: []string{student1ID},
	})
	if err != nil {
		t.Fatalf("expected s1 created successfully, got: %v", err)
	}
	if s1.Type != SessionTypeIndividual || s1.DurationMinutes != 75 {
		t.Errorf("unexpected s1 props: type=%s, duration=%d", s1.Type, s1.DurationMinutes)
	}

	// 2. Create CONSECUTIVE second session: 09:15 - 10:30 (NO mandatory break!)
	// Rule: Must be accepted because no breaks are required.
	s2, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "09:15",
		EndTime:    "10:30",
		StudentIDs: []string{student2ID},
	})
	if err != nil {
		t.Fatalf("consecutive session must be accepted without mandatory break, got error: %v", err)
	}
	if s2.StartTime != "09:15" || s2.EndTime != "10:30" {
		t.Errorf("unexpected s2 times: %s - %s", s2.StartTime, s2.EndTime)
	}

	// 3. Reject overlapping session: 09:00 - 10:00 (overlaps with s1 and s2)
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "09:00",
		EndTime:    "10:00",
		StudentIDs: []string{student1ID},
	})
	if !errors.Is(err, ErrSessionOverlap) {
		t.Errorf("expected ErrSessionOverlap, got: %v", err)
	}

	// 4. Reject session outside coach availability: 14:00 - 15:15
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "14:00",
		EndTime:    "15:15",
		StudentIDs: []string{student1ID},
	})
	if !errors.Is(err, ErrSessionOutsideAvailability) {
		t.Errorf("expected ErrSessionOutsideAvailability, got: %v", err)
	}

	// 5. Reject session with archived student
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "11:00",
		EndTime:    "12:15",
		StudentIDs: []string{archivedID},
	})
	if !errors.Is(err, ErrStudentNotFoundOrArchived) {
		t.Errorf("expected ErrStudentNotFoundOrArchived, got: %v", err)
	}

	// 6. Group session with multiple students (no arbitrary cap!)
	sGroup, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "11:00",
		EndTime:    "12:15",
		StudentIDs: []string{student1ID, student2ID},
	})
	if err != nil {
		t.Fatalf("failed creating group session: %v", err)
	}
	if sGroup.Type != SessionTypeGroup || len(sGroup.StudentIDs) != 2 {
		t.Errorf("expected group session with 2 students, got type=%s, count=%d", sGroup.Type, len(sGroup.StudentIDs))
	}

	// 7. Reject duplicate student in same session
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mié",
		StartTime:  "11:00",
		EndTime:    "12:15",
		StudentIDs: []string{student1ID, student1ID},
	})
	if !errors.Is(err, ErrDuplicateStudentInSession) {
		t.Errorf("expected ErrDuplicateStudentInSession, got: %v", err)
	}
}

func TestAvailabilityConflictWithExistingSessions(t *testing.T) {
	repo := newInMemoryRepo()
	studentID := primitive.NewObjectID().Hex()
	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			studentID: {ID: studentID, FullName: "Alumno Test", Archived: false},
		},
	}
	svc := NewService(repo, reader)
	ctx := context.Background()

	// Initial availability: Friday 08:00 - 18:00
	enabled := true
	_, err := svc.UpdateAvailabilityForDay(ctx, "Vie", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "18:00"},
		},
	})
	if err != nil {
		t.Fatalf("failed setting availability: %v", err)
	}

	// Active session at 15:00 - 16:15
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Vie",
		StartTime:  "15:00",
		EndTime:    "16:15",
		StudentIDs: []string{studentID},
	})
	if err != nil {
		t.Fatalf("failed creating session: %v", err)
	}

	// Now try to reduce availability to morning only (08:00 - 12:00)
	// Must be rejected with AvailabilityConflictError identifying the session!
	_, err = svc.UpdateAvailabilityForDay(ctx, "Vie", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "12:00"},
		},
	})

	var conflictErr *AvailabilityConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("expected AvailabilityConflictError, got: %v", err)
	}
	if len(conflictErr.Conflicts) != 1 {
		t.Errorf("expected 1 conflicting session, got %d", len(conflictErr.Conflicts))
	}
}

func TestCalendarExceptionsAndProjection(t *testing.T) {
	repo := newInMemoryRepo()
	studentID := primitive.NewObjectID().Hex()
	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			studentID: {ID: studentID, FullName: "Alumno Calendario", Archived: false},
		},
	}
	svc := NewService(repo, reader)
	ctx := context.Background()

	// Availability for Monday (Lun) and Thursday (Jue)
	enabled := true
	_, _ = svc.UpdateAvailabilityForDay(ctx, "Lun", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "12:00"},
		},
	})
	_, _ = svc.UpdateAvailabilityForDay(ctx, "Jue", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "10:00", EndTime: "14:00"},
		},
	})

	// Recurring Monday session starting 2026-11-02 (Monday)
	sess, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "08:00",
		EndTime:    "09:15",
		StartDate:  "2026-11-01",
		StudentIDs: []string{studentID},
	})
	if err != nil {
		t.Fatalf("failed creating recurring session: %v", err)
	}

	// 1. Project occurrences for week of 2026-11-02 to 2026-11-08
	occs, err := svc.GetOccurrences(ctx, SessionsFilter{
		From: "2026-11-02",
		To:   "2026-11-08",
	})
	if err != nil {
		t.Fatalf("failed getting occurrences: %v", err)
	}
	if len(occs) != 1 || occs[0].Date != "2026-11-02" {
		t.Fatalf("expected 1 occurrence on 2026-11-02, got: %d", len(occs))
	}

	// 2. Reschedule Monday occurrence (2026-11-02) to Thursday 2026-11-05 10:00 - 11:15
	newDate := "2026-11-05"
	newStart := "10:00"
	newEnd := "11:15"
	exc, err := svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
		Date:         "2026-11-02",
		Type:         ExceptionTypeRescheduled,
		NewDate:      &newDate,
		NewStartTime: &newStart,
		NewEndTime:   &newEnd,
	})
	if err != nil {
		t.Fatalf("failed creating reschedule exception: %v", err)
	}
	if exc.Type != ExceptionTypeRescheduled {
		t.Errorf("unexpected exception type: %s", exc.Type)
	}

	// 3. Project occurrences again:
	// - 2026-11-02 should have NO occurrence (it was moved away)
	// - 2026-11-05 should have the rescheduled occurrence!
	occs, err = svc.GetOccurrences(ctx, SessionsFilter{
		From: "2026-11-02",
		To:   "2026-11-08",
	})
	if err != nil {
		t.Fatalf("failed getting occurrences after exception: %v", err)
	}
	if len(occs) != 1 {
		t.Fatalf("expected exactly 1 rescheduled occurrence in week, got %d", len(occs))
	}
	if occs[0].Date != "2026-11-05" || !occs[0].IsException || occs[0].StartTime != "10:00" {
		t.Errorf("expected rescheduled occurrence on 2026-11-05 at 10:00, got: %+v", occs[0])
	}

	// 4. Verify that recurring session rule was NOT altered
	persistedSess, err := svc.GetSessionByID(ctx, sess.ID)
	if err != nil {
		t.Fatalf("failed getting original session: %v", err)
	}
	if persistedSess.Day != DayLun || persistedSess.StartTime != "08:00" {
		t.Errorf("original recurring session was improperly modified: %+v", persistedSess)
	}

	// 5. Delete exception: original occurrence is restored!
	err = svc.DeleteException(ctx, sess.ID, exc.ID)
	if err != nil {
		t.Fatalf("failed deleting exception: %v", err)
	}

	occs, err = svc.GetOccurrences(ctx, SessionsFilter{
		From: "2026-11-02",
		To:   "2026-11-08",
	})
	if err != nil {
		t.Fatalf("failed getting occurrences after restoring: %v", err)
	}
	if len(occs) != 1 || occs[0].Date != "2026-11-02" || occs[0].IsException {
		t.Errorf("expected original occurrence restored on 2026-11-02, got: %+v", occs)
	}
}

func TestSessionConflictRegressionRules(t *testing.T) {
	ctx := context.Background()
	repo := newInMemoryRepo()
	student1 := primitive.NewObjectID().Hex()
	student2 := primitive.NewObjectID().Hex()
	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			student1: {ID: student1, FullName: "Alumno Regresión 1"},
			student2: {ID: student2, FullName: "Alumno Regresión 2"},
		},
	}
	svc := NewService(repo, reader)

	// Setup coach availability: Monday 08:00 - 20:00, Tuesday 08:00 - 20:00
	enabled := true
	_, _ = svc.UpdateAvailabilityForDay(ctx, "Lun", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "20:00"},
		},
	})
	_, _ = svc.UpdateAvailabilityForDay(ctx, "Mar", UpdateAvailabilityRequest{
		Enabled: &enabled,
		Intervals: []AvailabilityInterval{
			{StartTime: "08:00", EndTime: "20:00"},
		},
	})

	// 1. Create base session: Monday 08:00 - 09:15 starting 2026-11-02
	sess1, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "08:00",
		EndTime:    "09:15",
		StartDate:  "2026-11-02",
		StudentIDs: []string{student1},
	})
	if err != nil {
		t.Fatalf("failed creating sess1: %v", err)
	}

	// Rule 1: Dos sesiones el mismo lunes y horario superpuesto -> Conflicto (ErrSessionOverlap)
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "08:30",
		EndTime:    "09:45",
		StartDate:  "2026-11-02",
		StudentIDs: []string{student2},
	})
	if !errors.Is(err, ErrSessionOverlap) {
		t.Errorf("Rule 1 failed: expected ErrSessionOverlap, got %v", err)
	}

	// Rule 2: Dos sesiones consecutivas (08:00-09:15 y 09:15-10:30) -> Permitidas
	sessConsecutive, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "09:15",
		EndTime:    "10:30",
		StartDate:  "2026-11-02",
		StudentIDs: []string{student2},
	})
	if err != nil {
		t.Errorf("Rule 2 failed: consecutive session rejected: %v", err)
	}
	if sessConsecutive == nil {
		t.Errorf("Rule 2 failed: expected session response")
	}

	// Rule 3: Dos sesiones del mismo día, pero con períodos de vigencia que nunca coinciden -> Permitidas
	// Create past session: Monday 14:00 - 15:15, startDate: 2026-09-01, endDate: 2026-09-30
	endPast := "2026-09-30"
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "14:00",
		EndTime:    "15:15",
		StartDate:  "2026-09-01",
		EndDate:    &endPast,
		StudentIDs: []string{student1},
	})
	if err != nil {
		t.Fatalf("failed creating past session: %v", err)
	}
	// New session in same slot 14:00 - 15:15, but starting in November 2026-11-01 -> No overlap in dates!
	sessNonOverlapDates, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "14:00",
		EndTime:    "15:15",
		StartDate:  "2026-11-01",
		StudentIDs: []string{student2},
	})
	if err != nil {
		t.Errorf("Rule 3 failed: sessions with disjoint date ranges must be allowed, got error: %v", err)
	}
	if sessNonOverlapDates == nil {
		t.Errorf("Rule 3 failed: expected valid session")
	}

	// Rule 4: Una sesión archivada y otra nueva en el mismo horario -> Permitida
	// Archive sessConsecutive (09:15 - 10:30)
	if err := svc.ArchiveSession(ctx, sessConsecutive.ID); err != nil {
		t.Fatalf("failed archiving session: %v", err)
	}
	// Re-create in 09:15 - 10:30 -> Must be allowed!
	sessAfterArchived, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "09:15",
		EndTime:    "10:30",
		StartDate:  "2026-11-02",
		StudentIDs: []string{student1},
	})
	if err != nil {
		t.Errorf("Rule 4 failed: session slot of archived session must be allowed, got %v", err)
	}
	if sessAfterArchived == nil {
		t.Errorf("Rule 4 failed: expected valid session")
	}

	// Rule 5: Una sesión activa existente dentro de su período de vigencia -> Conflicto
	// sess1 is active (08:00-09:15 from 2026-11-02 indefinite)
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "08:00",
		EndTime:    "09:15",
		StartDate:  "2026-11-16",
		StudentIDs: []string{student2},
	})
	if !errors.Is(err, ErrSessionOverlap) {
		t.Errorf("Rule 5 failed: expected ErrSessionOverlap with active session in validity period, got %v", err)
	}

	// Rule 6: Una ocurrencia cancelada mediante excepción -> No debe generar conflicto para esa fecha específica
	// Cancel sess1 on date 2026-11-23
	cancelDate := "2026-11-23"
	_, err = svc.CreateException(ctx, sess1.ID, CreateExceptionRequest{
		Date: cancelDate,
		Type: ExceptionTypeCancelled,
	})
	if err != nil {
		t.Fatalf("failed creating cancel exception: %v", err)
	}
	// A single-date session on 2026-11-23 08:00-09:15 -> Must be allowed because original occurrence was cancelled!
	sessOnCancelledDate, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "08:00",
		EndTime:    "09:15",
		StartDate:  cancelDate,
		EndDate:    &cancelDate,
		StudentIDs: []string{student2},
	})
	if err != nil {
		t.Errorf("Rule 6 failed: session on date with cancelled exception must be allowed, got: %v", err)
	}
	if sessOnCancelledDate == nil {
		t.Errorf("Rule 6 failed: expected valid session")
	}

	// Rule 7: Una ocurrencia reprogramada: el horario original queda libre en esa fecha y el nuevo horario valida conflictos
	// Reschedule sessAfterArchived (09:15-10:30) on 2026-11-02 away to Tuesday 2026-11-03 at 11:00-12:15
	reschedDate := "2026-11-02"
	destDate := "2026-11-03"
	destStart := "11:00"
	destEnd := "12:15"
	_, err = svc.CreateException(ctx, sessAfterArchived.ID, CreateExceptionRequest{
		Date:         reschedDate,
		Type:         ExceptionTypeRescheduled,
		NewDate:      &destDate,
		NewStartTime: &destStart,
		NewEndTime:   &destEnd,
	})
	if err != nil {
		t.Fatalf("failed rescheduling session: %v", err)
	}
	// The original slot on 2026-11-02 (09:15-10:30) should now be free on that single date:
	sessOnFreedReschedDate, err := svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Lun",
		StartTime:  "09:15",
		EndTime:    "10:30",
		StartDate:  reschedDate,
		EndDate:    &reschedDate,
		StudentIDs: []string{student2},
	})
	if err != nil {
		t.Errorf("Rule 7a failed: original slot freed by reschedule must be available, got: %v", err)
	}
	if sessOnFreedReschedDate == nil {
		t.Errorf("Rule 7a failed: expected valid session")
	}
	// And the new destination slot on 2026-11-03 at 11:00-12:15 must reject conflicts:
	_, err = svc.CreateSession(ctx, CreateSessionRequest{
		Day:        "Mar",
		StartTime:  "11:00",
		EndTime:    "12:15",
		StartDate:  destDate,
		EndDate:    &destDate,
		StudentIDs: []string{student1},
	})
	if !errors.Is(err, ErrSessionOverlap) {
		t.Errorf("Rule 7b failed: destination slot occupied by rescheduled session must reject conflict, got %v", err)
	}

	// Rule 8: La proyección semanal debe mostrar todas las sesiones que pueden generar conflictos para esa semana
	occs, err := svc.GetOccurrences(ctx, SessionsFilter{
		From: "2026-11-02",
		To:   "2026-11-08",
	})
	if err != nil {
		t.Fatalf("Rule 8 failed: error fetching occurrences: %v", err)
	}
	if len(occs) < 3 {
		t.Errorf("Rule 8 failed: expected at least 3 occurrences visible for this week, got %d", len(occs))
	}
	for _, occ := range occs {
		if occ.Date < "2026-11-02" || occ.Date > "2026-11-08" {
			t.Errorf("Rule 8 failed: occurrence %s falls outside requested week", occ.Date)
		}
	}
}

func TestSynchronizedOccurrencesAndConflictRules(t *testing.T) {
	ctx := context.Background()
	repo := newInMemoryRepo()
	student1 := primitive.NewObjectID().Hex()
	student2 := primitive.NewObjectID().Hex()
	reader := &mockStudentReader{
		students: map[string]StudentSummary{
			student1: {ID: student1, FullName: "Alumno Test 1"},
			student2: {ID: student2, FullName: "Alumno Test 2"},
		},
	}
	svc := NewService(repo, reader)

	enabled := true
	for _, day := range []DayOfWeek{"Lun", "Mar", "Mie", "Jue", "Vie", "Sab"} {
		_, _ = svc.UpdateAvailabilityForDay(ctx, string(day), UpdateAvailabilityRequest{
			Enabled: &enabled,
			Intervals: []AvailabilityInterval{
				{StartTime: "07:00", EndTime: "22:00"},
			},
		})
	}

	// -------------------------------------------------------------
	// EXACT REPORTED BUG REPRODUCTION:
	// A parent session on Monday 2026-11-02 (09:15-10:30) is archived.
	// An orphan exception exists on 2026-11-02 (09:20-10:35).
	// Attempting to book Monday 10:00-11:15 must NOT trigger a phantom conflict!
	// -------------------------------------------------------------
	t.Run("ReportedCase_ArchivedParentRescheduledException_NoPhantomConflict", func(t *testing.T) {
		archivedSessID := primitive.NewObjectID()
		repo.sessions[archivedSessID] = Session{
			ID:              archivedSessID,
			Day:             "Lun",
			StartTime:       "09:15",
			EndTime:         "10:30",
			DurationMinutes: 75,
			Type:            SessionTypeIndividual,
			StudentIDs:      []primitive.ObjectID{primitive.NewObjectID()},
			IsRecurring:     true,
			StartDate:       "2026-11-02",
			Status:          SessionStatusArchived, // ARCHIVED!
		}
		// Residual exception in DB pointing to archived session
		excID := primitive.NewObjectID()
		newDate := "2026-11-02"
		newDay := DayOfWeek("Lun")
		newStart := "09:20"
		newEnd := "10:35"
		repo.exceptions[excID] = ScheduleException{
			ID:           excID,
			SessionID:    archivedSessID,
			OriginalDate: "2026-11-02",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewDay:       &newDay,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		}

		// 1. Verify calendar projection does NOT show the archived session or its exception
		occs, err := svc.GetOccurrences(ctx, SessionsFilter{
			From: "2026-11-02",
			To:   "2026-11-08",
		})
		if err != nil {
			t.Fatalf("unexpected error fetching occurrences: %v", err)
		}
		for _, occ := range occs {
			if occ.BlockID == archivedSessID.Hex() || occ.StartTime == "09:20" {
				t.Fatalf("archived session exception should not appear in occurrences: %+v", occ)
			}
		}

		// 2. Verify ListExceptions does NOT return the exception of archived session
		excs, err := svc.ListExceptions(ctx, ExceptionsFilter{
			From: "2026-11-02",
			To:   "2026-11-08",
		})
		if err != nil {
			t.Fatalf("unexpected error listing exceptions: %v", err)
		}
		for _, e := range excs {
			if e.SessionID == archivedSessID.Hex() {
				t.Fatalf("archived session exception should not appear in ListExceptions: %+v", e)
			}
		}

		// 3. Attempting to book Monday 10:00-11:15 MUST SUCCEED (no phantom conflict!)
		newSess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Lun",
			StartTime:  "10:00",
			EndTime:    "11:15",
			StartDate:  "2026-11-02",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("expected booking 10:00-11:15 to succeed, got phantom conflict: %v", err)
		}
		if newSess == nil || newSess.StartTime != "10:00" {
			t.Fatalf("expected created session with StartTime 10:00, got: %+v", newSess)
		}

		// Cleanup for subsequent isolated subtests
		_ = svc.ArchiveSession(ctx, newSess.ID)
		delete(repo.sessions, archivedSessID)
		delete(repo.exceptions, excID)
	})

	// -------------------------------------------------------------
	// SCENARIO 1: Rescheduled session into visible range -> appears in calendar
	// -------------------------------------------------------------
	t.Run("Scenario1_RescheduledSessionIntoVisibleRange", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Lun",
			StartTime:  "08:00",
			EndTime:    "09:15",
			StartDate:  "2026-11-02",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Reschedule same day to 09:20 - 10:35
		newDate := "2026-11-02"
		newStart := "09:20"
		newEnd := "10:35"
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date:         "2026-11-02",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		})
		if err != nil {
			t.Fatalf("failed rescheduling: %v", err)
		}

		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		var foundResched, foundOriginal bool
		for _, occ := range occs {
			if occ.Date == "2026-11-02" && occ.StartTime == "09:20" && occ.EndTime == "10:35" && occ.IsException {
				foundResched = true
			}
			if occ.Date == "2026-11-02" && occ.StartTime == "08:00" {
				foundOriginal = true
			}
		}
		if !foundResched {
			t.Errorf("Scenario 1 failed: rescheduled occurrence 09:20-10:35 not found in calendar")
		}
		if foundOriginal {
			t.Errorf("Scenario 1 failed: original slot 08:00 should not be projected after reschedule")
		}
	})

	// -------------------------------------------------------------
	// SCENARIO 2: Original date outside range, new date inside -> appears in calendar
	// -------------------------------------------------------------
	t.Run("Scenario2_OriginalDateOutside_NewDateInside", func(t *testing.T) {
		// Session starts in October
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Vie",
			StartTime:  "14:00",
			EndTime:    "15:15",
			StartDate:  "2026-10-30",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Rescheduled from 2026-10-30 (outside Nov week) to 2026-11-04 (Mie, inside Nov week)
		newDate := "2026-11-04"
		newStart := "14:00"
		newEnd := "15:15"
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date:         "2026-10-30",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		})
		if err != nil {
			t.Fatalf("failed creating exception: %v", err)
		}

		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		found := false
		for _, occ := range occs {
			if occ.Date == "2026-11-04" && occ.StartTime == "14:00" && occ.IsException {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Scenario 2 failed: session rescheduled from outside range into 2026-11-04 must appear")
		}
	})

	// -------------------------------------------------------------
	// SCENARIO 3: Original date inside range, new date outside -> does not appear on original date
	// -------------------------------------------------------------
	t.Run("Scenario3_OriginalDateInside_NewDateOutside", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Mie",
			StartTime:  "16:00",
			EndTime:    "17:15",
			StartDate:  "2026-11-04",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Reschedule from 2026-11-04 to 2026-11-20 (outside current week)
		newDate := "2026-11-20"
		newStart := "16:00"
		newEnd := "17:15"
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date:         "2026-11-04",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		})
		if err != nil {
			t.Fatalf("failed creating exception: %v", err)
		}

		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		for _, occ := range occs {
			if occ.Date == "2026-11-04" && occ.StartTime == "16:00" {
				t.Fatalf("Scenario 3 failed: session moved outside week should not appear on original date 2026-11-04")
			}
		}
	})

	// -------------------------------------------------------------
	// SCENARIO 4: Cancelled occurrence -> does not appear nor generate conflict
	// -------------------------------------------------------------
	t.Run("Scenario4_CancelledOccurrence_NoAppear_NoConflict", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Jue",
			StartTime:  "10:00",
			EndTime:    "11:15",
			StartDate:  "2026-11-05",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Cancel on 2026-11-05
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date: "2026-11-05",
			Type: ExceptionTypeCancelled,
		})
		if err != nil {
			t.Fatalf("failed creating cancellation: %v", err)
		}

		// Verify it does not appear on 2026-11-05
		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		for _, occ := range occs {
			if occ.Date == "2026-11-05" && occ.StartTime == "10:00" {
				t.Fatalf("Scenario 4 failed: cancelled session still appears in occurrences")
			}
		}

		// Verify booking on that specific cancelled date does NOT conflict
		singleDate := "2026-11-05"
		replacementSess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Jue",
			StartTime:  "10:00",
			EndTime:    "11:15",
			StartDate:  singleDate,
			EndDate:    &singleDate,
			StudentIDs: []string{student2},
		})
		if err != nil {
			t.Fatalf("Scenario 4 failed: cancelled slot should be bookable on that single date, got: %v", err)
		}
		_ = svc.ArchiveSession(ctx, replacementSess.ID)
	})

	// -------------------------------------------------------------
	// SCENARIO 5: Rescheduled occurrence -> original slot freed (no conflict on original slot)
	// -------------------------------------------------------------
	t.Run("Scenario5_RescheduledOccurrence_OriginalSlotFreed", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Vie",
			StartTime:  "11:00",
			EndTime:    "12:15",
			StartDate:  "2026-11-06",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Reschedule to 15:00-16:15
		newDate := "2026-11-06"
		newStart := "15:00"
		newEnd := "16:15"
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date:         "2026-11-06",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		})
		if err != nil {
			t.Fatalf("failed rescheduling: %v", err)
		}

		// Single date booking on original slot (11:00-12:15) must be permitted
		singleDate := "2026-11-06"
		freedSlotSess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Vie",
			StartTime:  "11:00",
			EndTime:    "12:15",
			StartDate:  singleDate,
			EndDate:    &singleDate,
			StudentIDs: []string{student2},
		})
		if err != nil {
			t.Fatalf("Scenario 5 failed: original slot freed by reschedule should be bookable, got: %v", err)
		}
		_ = svc.ArchiveSession(ctx, freedSlotSess.ID)
	})

	// -------------------------------------------------------------
	// SCENARIO 6: Attempting to occupy new rescheduled slot -> conflict
	// -------------------------------------------------------------
	t.Run("Scenario6_OccupyingNewRescheduledSlot_GeneratesConflict", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Sab",
			StartTime:  "09:00",
			EndTime:    "10:15",
			StartDate:  "2026-11-07",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		// Reschedule to 12:00 - 13:15
		newDate := "2026-11-07"
		newStart := "12:00"
		newEnd := "13:15"
		_, err = svc.CreateException(ctx, sess.ID, CreateExceptionRequest{
			Date:         "2026-11-07",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		})
		if err != nil {
			t.Fatalf("failed rescheduling: %v", err)
		}

		// Attempting to occupy 12:30 - 13:45 on 2026-11-07 must conflict
		_, err = svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Sab",
			StartTime:  "12:30",
			EndTime:    "13:45",
			StartDate:  "2026-11-07",
			StudentIDs: []string{student2},
		})
		if !errors.Is(err, ErrSessionOverlap) {
			t.Fatalf("Scenario 6 failed: expected ErrSessionOverlap when booking overlapping rescheduled slot, got: %v", err)
		}
	})

	// -------------------------------------------------------------
	// SCENARIO 7: Archived parent session -> no active occurrences or phantom conflicts
	// -------------------------------------------------------------
	t.Run("Scenario7_ArchivedParentSession_NoActiveOccurrencesOrConflicts", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Mar",
			StartTime:  "09:00",
			EndTime:    "10:15",
			StartDate:  "2026-11-03",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		// Archive it immediately
		err = svc.ArchiveSession(ctx, sess.ID)
		if err != nil {
			t.Fatalf("failed archiving session: %v", err)
		}

		// Calendar check
		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		for _, occ := range occs {
			if occ.BlockID == sess.ID {
				t.Fatalf("Scenario 7 failed: archived session should not appear in occurrences")
			}
		}

		// Conflict check: slot 09:00-10:15 should be fully free
		newSess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Mar",
			StartTime:  "09:00",
			EndTime:    "10:15",
			StartDate:  "2026-11-03",
			StudentIDs: []string{student2},
		})
		if err != nil {
			t.Fatalf("Scenario 7 failed: slot of archived session must be bookable, got %v", err)
		}
		_ = svc.ArchiveSession(ctx, newSess.ID)
	})

	// -------------------------------------------------------------
	// SCENARIO 8: Exception with missing/archived session -> ignored, no phantom conflict
	// -------------------------------------------------------------
	t.Run("Scenario8_ExceptionWithMissingSession_Ignored", func(t *testing.T) {
		missingSessionID := primitive.NewObjectID()
		orphanExcID := primitive.NewObjectID()
		newDate := "2026-11-04"
		newDay := DayOfWeek("Mie")
		newStart := "18:00"
		newEnd := "19:15"
		repo.exceptions[orphanExcID] = ScheduleException{
			ID:           orphanExcID,
			SessionID:    missingSessionID,
			OriginalDate: "2026-11-04",
			Type:         ExceptionTypeRescheduled,
			NewDate:      &newDate,
			NewDay:       &newDay,
			NewStartTime: &newStart,
			NewEndTime:   &newEnd,
		}

		// Must not produce occurrences
		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		for _, occ := range occs {
			if occ.StartTime == "18:00" && occ.Date == "2026-11-04" {
				t.Fatalf("Scenario 8 failed: orphan exception generated occurrence")
			}
		}

		// Must not produce conflict
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Mie",
			StartTime:  "18:00",
			EndTime:    "19:15",
			StartDate:  "2026-11-04",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("Scenario 8 failed: orphan exception caused phantom conflict: %v", err)
		}
		_ = svc.ArchiveSession(ctx, sess.ID)
		delete(repo.exceptions, orphanExcID)
	})

	// -------------------------------------------------------------
	// SCENARIO 9: Regular week without exceptions -> current projection maintained
	// -------------------------------------------------------------
	t.Run("Scenario9_RegularWeekWithoutExceptions_CurrentProjectionMaintained", func(t *testing.T) {
		sess, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Jue",
			StartTime:  "08:00",
			EndTime:    "09:15",
			StartDate:  "2026-11-01",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating session: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess.ID) }()

		occs, err := svc.GetOccurrences(ctx, SessionsFilter{From: "2026-11-02", To: "2026-11-08"})
		if err != nil {
			t.Fatalf("failed getting occurrences: %v", err)
		}
		var matched *SessionOccurrenceResponse
		for i := range occs {
			if occs[i].BlockID == sess.ID && occs[i].Date == "2026-11-05" {
				matched = &occs[i]
				break
			}
		}
		if matched == nil {
			t.Fatalf("Scenario 9 failed: regular session not found on 2026-11-05")
		}
		if matched.IsException || matched.StartTime != "08:00" || matched.EndTime != "09:15" {
			t.Fatalf("Scenario 9 failed: incorrect projection data: %+v", matched)
		}
	})

	// -------------------------------------------------------------
	// SCENARIO 10: Consecutive sessions -> allowed
	// -------------------------------------------------------------
	t.Run("Scenario10_ConsecutiveSessions_Allowed", func(t *testing.T) {
		sess1, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Vie",
			StartTime:  "08:00",
			EndTime:    "09:15",
			StartDate:  "2026-11-06",
			StudentIDs: []string{student1},
		})
		if err != nil {
			t.Fatalf("failed creating sess1: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess1.ID) }()

		// Consecutive: 09:15 - 10:30 (starts exactly when sess1 ends)
		sess2, err := svc.CreateSession(ctx, CreateSessionRequest{
			Day:        "Vie",
			StartTime:  "09:15",
			EndTime:    "10:30",
			StartDate:  "2026-11-06",
			StudentIDs: []string{student2},
		})
		if err != nil {
			t.Fatalf("Scenario 10 failed: consecutive session (09:15-10:30) rejected: %v", err)
		}
		defer func() { _ = svc.ArchiveSession(ctx, sess2.ID) }()
	})
}
