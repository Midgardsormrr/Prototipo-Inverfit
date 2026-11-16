package schedule

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
	_ "time/tzdata" // embed IANA timezone database for zero-dependency timezone support on Windows

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AvailabilityConflictError is returned when an availability reduction would strand active sessions.
type AvailabilityConflictError struct {
	Message   string
	Conflicts []SessionResponse
}

func (e *AvailabilityConflictError) Error() string {
	return e.Message
}

// Service defines the business logic contract for the Schedule module.
type Service interface {
	// Availability
	GetAvailability(ctx context.Context) ([]DayAvailabilityResponse, error)
	GetAvailabilityByDay(ctx context.Context, dayStr string) (*DayAvailabilityResponse, error)
	UpdateAvailabilityForDay(ctx context.Context, dayStr string, req UpdateAvailabilityRequest) (*DayAvailabilityResponse, error)
	DeleteAvailabilityForDay(ctx context.Context, dayStr string) error

	// Sessions
	CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error)
	GetSessionByID(ctx context.Context, idStr string) (*SessionResponse, error)
	UpdateSession(ctx context.Context, idStr string, req UpdateSessionRequest) (*SessionResponse, error)
	ArchiveSession(ctx context.Context, idStr string) error
	GetOccurrences(ctx context.Context, filter SessionsFilter) ([]SessionOccurrenceResponse, error)

	// Exceptions
	CreateException(ctx context.Context, sessionIDStr string, req CreateExceptionRequest) (*ExceptionResponse, error)
	ListExceptions(ctx context.Context, filter ExceptionsFilter) ([]ExceptionResponse, error)
	DeleteException(ctx context.Context, sessionIDStr, exceptionIDStr string) error
}

type service struct {
	repo          Repository
	studentReader StudentReader
	loc           *time.Location
}

// NewService instantiates the schedule business service with injected dependencies.
func NewService(repo Repository, studentReader StudentReader) Service {
	loc, err := time.LoadLocation(DefaultTimezone)
	if err != nil {
		loc = time.FixedZone("America/Santiago", -3*3600) // fallback UTC-3
	}

	return &service{
		repo:          repo,
		studentReader: studentReader,
		loc:           loc,
	}
}

// ==========================================
// AVAILABILITY METHODS
// ==========================================

func (s *service) GetAvailability(ctx context.Context) ([]DayAvailabilityResponse, error) {
	stored, err := s.repo.GetAllAvailability(ctx)
	if err != nil {
		return nil, err
	}

	storedMap := make(map[DayOfWeek]DayAvailability)
	for _, a := range stored {
		storedMap[a.Day] = a
	}

	// Always return all 7 canonical days in order
	res := make([]DayAvailabilityResponse, len(CanonicalDays))
	for i, day := range CanonicalDays {
		if a, ok := storedMap[day]; ok {
			intervals := a.Intervals
			if intervals == nil {
				intervals = make([]AvailabilityInterval, 0)
			}
			res[i] = DayAvailabilityResponse{
				Day:       a.Day,
				Enabled:   a.Enabled,
				Intervals: intervals,
				Timezone:  a.Timezone,
			}
		} else {
			res[i] = DayAvailabilityResponse{
				Day:       day,
				Enabled:   false,
				Intervals: make([]AvailabilityInterval, 0),
				Timezone:  DefaultTimezone,
			}
		}
	}

	return res, nil
}

func (s *service) GetAvailabilityByDay(ctx context.Context, dayStr string) (*DayAvailabilityResponse, error) {
	day, err := NormalizeDayOfWeek(dayStr)
	if err != nil {
		return nil, err
	}

	stored, err := s.repo.GetAvailabilityByDay(ctx, day)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return &DayAvailabilityResponse{
			Day:       day,
			Enabled:   false,
			Intervals: make([]AvailabilityInterval, 0),
			Timezone:  DefaultTimezone,
		}, nil
	}

	intervals := stored.Intervals
	if intervals == nil {
		intervals = make([]AvailabilityInterval, 0)
	}

	return &DayAvailabilityResponse{
		Day:       stored.Day,
		Enabled:   stored.Enabled,
		Intervals: intervals,
		Timezone:  stored.Timezone,
	}, nil
}

func (s *service) UpdateAvailabilityForDay(ctx context.Context, dayStr string, req UpdateAvailabilityRequest) (*DayAvailabilityResponse, error) {
	day, err := NormalizeDayOfWeek(dayStr)
	if err != nil {
		return nil, err
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	// Validate intervals if enabled
	intervals := req.Intervals
	if intervals == nil {
		intervals = []AvailabilityInterval{}
	}

	if enabled && len(intervals) > 0 {
		for i := range intervals {
			if intervals[i].ID == "" {
				intervals[i].ID = fmt.Sprintf("int-%s-%d", strings.ToLower(string(day)), time.Now().UnixNano()+int64(i))
			}
			startM, err := TimeToMinutes(intervals[i].StartTime)
			if err != nil {
				return nil, fmt.Errorf("%w: hora de inicio '%s'", ErrInvalidTimeFormat, intervals[i].StartTime)
			}
			endM, err := TimeToMinutes(intervals[i].EndTime)
			if err != nil {
				return nil, fmt.Errorf("%w: hora de término '%s'", ErrInvalidTimeFormat, intervals[i].EndTime)
			}
			if endM <= startM {
				return nil, ErrEndTimeBeforeStartTime
			}

			// Operating hours check (06:30 - 22:00)
			opStart, _ := TimeToMinutes(OperatingStart)
			opEnd, _ := TimeToMinutes(OperatingEnd)
			if startM < opStart || endM > opEnd {
				return nil, ErrOutsideOperatingHours
			}
		}

		// Check for overlaps between intervals of the same day
		for i := 0; i < len(intervals); i++ {
			startA, _ := TimeToMinutes(intervals[i].StartTime)
			endA, _ := TimeToMinutes(intervals[i].EndTime)
			for j := i + 1; j < len(intervals); j++ {
				startB, _ := TimeToMinutes(intervals[j].StartTime)
				endB, _ := TimeToMinutes(intervals[j].EndTime)
				if IntervalsOverlap(startA, endA, startB, endB) {
					return nil, ErrOverlappingIntervals
				}
			}
		}
	}

	// Check for conflicts with existing active sessions on this day
	activeSessions, err := s.repo.ListActiveSessions(ctx, &day, nil)
	if err != nil {
		return nil, err
	}

	conflictingSessions := make([]SessionResponse, 0)
	for _, session := range activeSessions {
		sStart, _ := TimeToMinutes(session.StartTime)
		sEnd, _ := TimeToMinutes(session.EndTime)

		if !enabled || len(intervals) == 0 {
			// All sessions become stranded
			conflictingSessions = append(conflictingSessions, ToSessionResponse(&session, nil))
			continue
		}

		// Must fit within at least one new interval
		fits := false
		for _, interval := range intervals {
			iStart, _ := TimeToMinutes(interval.StartTime)
			iEnd, _ := TimeToMinutes(interval.EndTime)
			if FitsInInterval(sStart, sEnd, iStart, iEnd) {
				fits = true
				break
			}
		}
		if !fits {
			conflictingSessions = append(conflictingSessions, ToSessionResponse(&session, nil))
		}
	}

	if len(conflictingSessions) > 0 {
		return nil, &AvailabilityConflictError{
			Message: fmt.Sprintf(
				"La nueva disponibilidad deja %d sesión(es) existente(s) fuera de horario. Resuelve los conflictos manualmente antes de guardar.",
				len(conflictingSessions),
			),
			Conflicts: conflictingSessions,
		}
	}

	// Persist availability
	avail := &DayAvailability{
		Day:       day,
		Enabled:   enabled,
		Intervals: intervals,
		Timezone:  DefaultTimezone,
	}

	if err := s.repo.UpsertAvailability(ctx, avail); err != nil {
		return nil, err
	}

	return &DayAvailabilityResponse{
		Day:       avail.Day,
		Enabled:   avail.Enabled,
		Intervals: avail.Intervals,
		Timezone:  avail.Timezone,
	}, nil
}

func (s *service) DeleteAvailabilityForDay(ctx context.Context, dayStr string) error {
	day, err := NormalizeDayOfWeek(dayStr)
	if err != nil {
		return err
	}

	// Check if active sessions exist
	activeSessions, err := s.repo.ListActiveSessions(ctx, &day, nil)
	if err != nil {
		return err
	}
	if len(activeSessions) > 0 {
		conflicts := make([]SessionResponse, 0, len(activeSessions))
		for _, s := range activeSessions {
			conflicts = append(conflicts, ToSessionResponse(&s, nil))
		}
		return &AvailabilityConflictError{
			Message: fmt.Sprintf(
				"No se puede eliminar la disponibilidad del %s porque existen %d sesión(es) activa(s) asignadas.",
				day,
				len(conflicts),
			),
			Conflicts: conflicts,
		}
	}

	return s.repo.DeleteAvailabilityByDay(ctx, day)
}

// ==========================================
// SESSIONS METHODS
// ==========================================

func (s *service) CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error) {
	day, err := NormalizeDayOfWeek(req.Day)
	if err != nil {
		return nil, err
	}

	startM, err := TimeToMinutes(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("%w: hora de inicio '%s'", ErrInvalidTimeFormat, req.StartTime)
	}
	endM, err := TimeToMinutes(req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("%w: hora de término '%s'", ErrInvalidTimeFormat, req.EndTime)
	}
	if endM <= startM {
		return nil, ErrEndTimeBeforeStartTime
	}
	duration := endM - startM
	if duration < 15 {
		return nil, ErrSessionDurationTooShort
	}

	// Operating hours
	opStart, _ := TimeToMinutes(OperatingStart)
	opEnd, _ := TimeToMinutes(OperatingEnd)
	if startM < opStart || endM > opEnd {
		return nil, ErrOutsideOperatingHours
	}

	// Validate students (no duplicates, must exist and not be archived)
	studentOIDs, summaries, err := s.validateStudents(ctx, req.StudentIDs)
	if err != nil {
		return nil, err
	}

	// Validate coach availability for that day
	dayAvail, err := s.repo.GetAvailabilityByDay(ctx, day)
	if err != nil {
		return nil, err
	}
	if dayAvail == nil || !dayAvail.Enabled || len(dayAvail.Intervals) == 0 {
		return nil, fmt.Errorf("%w para el día %s", ErrSessionOutsideAvailability, day)
	}
	fitsInAvailability := false
	for _, interval := range dayAvail.Intervals {
		iStart, _ := TimeToMinutes(interval.StartTime)
		iEnd, _ := TimeToMinutes(interval.EndTime)
		if FitsInInterval(startM, endM, iStart, iEnd) {
			fitsInAvailability = true
			break
		}
	}
	if !fitsInAvailability {
		return nil, fmt.Errorf("%w (%s — %s no está dentro de los intervalos configurados para el %s)", ErrSessionOutsideAvailability, req.StartTime, req.EndTime, day)
	}

	// Defaults and validation for startDate and endDate
	startDate := strings.TrimSpace(req.StartDate)
	if startDate == "" {
		now := time.Now().In(s.loc)
		weekday := int(now.Weekday())
		mondayOffset := 1 - weekday
		if weekday == 0 {
			mondayOffset = -6
		}
		monday := now.AddDate(0, 0, mondayOffset)
		targetWeekday, err := DayOfWeekToTimeWeekday(day)
		if err == nil {
			dayOffset := int(targetWeekday) - 1
			if targetWeekday == time.Sunday {
				dayOffset = 6
			}
			sessionDate := monday.AddDate(0, 0, dayOffset)
			startDate = sessionDate.Format("2006-01-02")
		} else {
			startDate = now.Format("2006-01-02")
		}
	} else if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, ErrInvalidDateFormat
	}

	var endDate *string
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		cleaned := strings.TrimSpace(*req.EndDate)
		if _, err := time.Parse("2006-01-02", cleaned); err != nil {
			return nil, ErrInvalidDateFormat
		}
		if cleaned < startDate {
			return nil, ErrInvalidDateRange
		}
		endDate = &cleaned
	}

	// Overlap check with other active sessions and exceptions considering date validity and recurrence
	if err := s.checkSessionConflicts(ctx, primitive.NilObjectID, day, startM, endM, startDate, endDate); err != nil {
		return nil, err
	}

	sessionType := SessionTypeIndividual
	if len(studentOIDs) > 1 {
		sessionType = SessionTypeGroup
	}

	isRecurring := true
	if req.IsRecurring != nil {
		isRecurring = *req.IsRecurring
	}

	session := &Session{
		ID:              primitive.NewObjectID(),
		Day:             day,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: duration,
		Type:            sessionType,
		StudentIDs:      studentOIDs,
		IsRecurring:     isRecurring,
		StartDate:       startDate,
		EndDate:         endDate,
		Status:          SessionStatusActive,
		Notes:           strings.TrimSpace(req.Notes),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	resp := ToSessionResponse(session, summaries)
	return &resp, nil
}

func (s *service) GetSessionByID(ctx context.Context, idStr string) (*SessionResponse, error) {
	oid, err := ParseHexID(idStr)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	session, err := s.repo.GetSessionByID(ctx, oid)
	if err != nil {
		return nil, err
	}

	var summaries []StudentSummary
	if s.studentReader != nil && len(session.StudentIDs) > 0 {
		strIDs := make([]string, len(session.StudentIDs))
		for i, id := range session.StudentIDs {
			strIDs[i] = id.Hex()
		}
		summaries, _ = s.studentReader.GetStudentsSummary(ctx, strIDs)
	}

	resp := ToSessionResponse(session, summaries)
	return &resp, nil
}

func (s *service) UpdateSession(ctx context.Context, idStr string, req UpdateSessionRequest) (*SessionResponse, error) {
	oid, err := ParseHexID(idStr)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	current, err := s.repo.GetSessionByID(ctx, oid)
	if err != nil {
		return nil, err
	}

	newDay := current.Day
	if req.Day != nil && *req.Day != "" {
		d, err := NormalizeDayOfWeek(*req.Day)
		if err != nil {
			return nil, err
		}
		newDay = d
	}

	newStart := current.StartTime
	if req.StartTime != nil && *req.StartTime != "" {
		newStart = *req.StartTime
	}

	newEnd := current.EndTime
	if req.EndTime != nil && *req.EndTime != "" {
		newEnd = *req.EndTime
	}

	startM, err := TimeToMinutes(newStart)
	if err != nil {
		return nil, fmt.Errorf("%w: hora de inicio '%s'", ErrInvalidTimeFormat, newStart)
	}
	endM, err := TimeToMinutes(newEnd)
	if err != nil {
		return nil, fmt.Errorf("%w: hora de término '%s'", ErrInvalidTimeFormat, newEnd)
	}
	if endM <= startM {
		return nil, ErrEndTimeBeforeStartTime
	}
	duration := endM - startM
	if duration < 15 {
		return nil, ErrSessionDurationTooShort
	}

	// Operating hours
	opStart, _ := TimeToMinutes(OperatingStart)
	opEnd, _ := TimeToMinutes(OperatingEnd)
	if startM < opStart || endM > opEnd {
		return nil, ErrOutsideOperatingHours
	}

	// Check coach availability for newDay
	dayAvail, err := s.repo.GetAvailabilityByDay(ctx, newDay)
	if err != nil {
		return nil, err
	}
	if dayAvail == nil || !dayAvail.Enabled || len(dayAvail.Intervals) == 0 {
		return nil, fmt.Errorf("%w para el día %s", ErrSessionOutsideAvailability, newDay)
	}
	fits := false
	for _, interval := range dayAvail.Intervals {
		iStart, _ := TimeToMinutes(interval.StartTime)
		iEnd, _ := TimeToMinutes(interval.EndTime)
		if FitsInInterval(startM, endM, iStart, iEnd) {
			fits = true
			break
		}
	}
	if !fits {
		return nil, fmt.Errorf("%w (%s — %s no está dentro de la disponibilidad configurada para el %s)", ErrSessionOutsideAvailability, newStart, newEnd, newDay)
	}

	startDate := current.StartDate
	if req.StartDate != nil && strings.TrimSpace(*req.StartDate) != "" {
		d := strings.TrimSpace(*req.StartDate)
		if _, err := time.Parse("2006-01-02", d); err != nil {
			return nil, ErrInvalidDateFormat
		}
		startDate = d
	}

	endDate := current.EndDate
	if req.EndDate != nil {
		if strings.TrimSpace(*req.EndDate) == "" {
			endDate = nil
		} else {
			d := strings.TrimSpace(*req.EndDate)
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return nil, ErrInvalidDateFormat
			}
			if d < startDate {
				return nil, ErrInvalidDateRange
			}
			endDate = &d
		}
	}

	// Check overlap with other active sessions and exceptions considering date validity and recurrence
	if err := s.checkSessionConflicts(ctx, current.ID, newDay, startM, endM, startDate, endDate); err != nil {
		return nil, err
	}

	// Validate students if provided
	studentOIDs := current.StudentIDs
	var summaries []StudentSummary
	if req.StudentIDs != nil {
		oids, sums, err := s.validateStudents(ctx, *req.StudentIDs)
		if err != nil {
			return nil, err
		}
		studentOIDs = oids
		summaries = sums
	} else if s.studentReader != nil && len(current.StudentIDs) > 0 {
		strIDs := make([]string, len(current.StudentIDs))
		for i, id := range current.StudentIDs {
			strIDs[i] = id.Hex()
		}
		summaries, _ = s.studentReader.GetStudentsSummary(ctx, strIDs)
	}

	sessionType := SessionTypeIndividual
	if len(studentOIDs) > 1 {
		sessionType = SessionTypeGroup
	}

	isRecurring := current.IsRecurring
	if req.IsRecurring != nil {
		isRecurring = *req.IsRecurring
	}

	notes := current.Notes
	if req.Notes != nil {
		notes = strings.TrimSpace(*req.Notes)
	}

	status := current.Status
	if req.Status != nil && *req.Status != "" {
		status = *req.Status
	}

	current.Day = newDay
	current.StartTime = newStart
	current.EndTime = newEnd
	current.DurationMinutes = duration
	current.Type = sessionType
	current.StudentIDs = studentOIDs
	current.IsRecurring = isRecurring
	current.StartDate = startDate
	current.EndDate = endDate
	current.Status = status
	current.Notes = notes

	if err := s.repo.UpdateSession(ctx, current); err != nil {
		return nil, err
	}

	resp := ToSessionResponse(current, summaries)
	return &resp, nil
}

func (s *service) ArchiveSession(ctx context.Context, idStr string) error {
	oid, err := ParseHexID(idStr)
	if err != nil {
		return ErrSessionNotFound
	}
	return s.repo.ArchiveSession(ctx, oid)
}

// ==========================================
// CALENDAR OCCURRENCES PROJECTION
// ==========================================

func (s *service) GetOccurrences(ctx context.Context, filter SessionsFilter) ([]SessionOccurrenceResponse, error) {
	from := strings.TrimSpace(filter.From)
	to := strings.TrimSpace(filter.To)

	// Default to current week Monday through Sunday if not specified
	now := time.Now().In(s.loc)
	if from == "" || to == "" {
		weekday := int(now.Weekday())
		mondayOffset := 1 - weekday
		if weekday == 0 {
			mondayOffset = -6
		}
		monday := now.AddDate(0, 0, mondayOffset)
		sunday := monday.AddDate(0, 0, 6)
		from = monday.Format("2006-01-02")
		to = sunday.Format("2006-01-02")
	}

	fromDate, err := time.ParseInLocation("2006-01-02", from, s.loc)
	if err != nil {
		return nil, fmt.Errorf("%w: 'from' debe ser YYYY-MM-DD", ErrInvalidDateFormat)
	}
	toDate, err := time.ParseInLocation("2006-01-02", to, s.loc)
	if err != nil {
		return nil, fmt.Errorf("%w: 'to' debe ser YYYY-MM-DD", ErrInvalidDateFormat)
	}
	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}

	var studentOID *primitive.ObjectID
	if filter.StudentID != "" {
		oid, err := ParseHexID(filter.StudentID)
		if err == nil {
			studentOID = &oid
		}
	}

	// Fetch all active recurring sessions
	sessions, err := s.repo.ListActiveSessions(ctx, nil, studentOID)
	if err != nil {
		return nil, err
	}

	// Fetch exceptions in period
	exceptions, err := s.repo.ListExceptionsByPeriod(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Build map of exceptions by session_id + original_date
	exceptionsMap := make(map[string]ScheduleException)
	rescheduledMovedToDate := make(map[string][]ScheduleException)
	for _, e := range exceptions {
		key := fmt.Sprintf("%s_%s", e.SessionID.Hex(), e.OriginalDate)
		exceptionsMap[key] = e
		if e.Type == ExceptionTypeRescheduled && e.NewDate != nil {
			rescheduledMovedToDate[*e.NewDate] = append(rescheduledMovedToDate[*e.NewDate], e)
		}
	}

	// Collect all student IDs to fetch summaries in a single batch
	studentIDSet := make(map[string]struct{})
	for _, s := range sessions {
		for _, sid := range s.StudentIDs {
			studentIDSet[sid.Hex()] = struct{}{}
		}
	}
	var allStudentIDs []string
	for sid := range studentIDSet {
		allStudentIDs = append(allStudentIDs, sid)
	}

	studentSummaryMap := make(map[string]StudentSummary)
	if s.studentReader != nil && len(allStudentIDs) > 0 {
		summaries, _ := s.studentReader.GetStudentsSummary(ctx, allStudentIDs)
		for _, sum := range summaries {
			studentSummaryMap[sum.ID] = sum
		}
	}

	occurrences := make([]SessionOccurrenceResponse, 0)

	// Iterate day by day from fromDate to toDate
	currentDate := fromDate
	for !currentDate.After(toDate) {
		dateStr := currentDate.Format("2006-01-02")
		dayOfWeek := TimeWeekdayToDayOfWeek(currentDate.Weekday())

		// 1. Process regular recurring sessions matching dayOfWeek
		for _, sess := range sessions {
			if sess.Day != dayOfWeek {
				continue
			}
			// Check date range of session
			if sess.StartDate > dateStr {
				continue
			}
			if sess.EndDate != nil && *sess.EndDate < dateStr {
				continue
			}

			// Check exceptions for this specific occurrence
			excKey := fmt.Sprintf("%s_%s", sess.ID.Hex(), dateStr)
			exc, hasExc := exceptionsMap[excKey]
			if hasExc {
				if exc.Type == ExceptionTypeCancelled || exc.Type == ExceptionTypeRescheduled {
					// Cancelado o reprogramado puntualmente fuera de este slot original.
					// Si fue reprogramado (incluso dentro del mismo día), su nueva ocurrencia
					// se proyecta abajo en movedExceptions con sus nuevos horarios.
					continue
				}
			}

			studentIDs := make([]string, 0, len(sess.StudentIDs))
			studentSummaries := make([]StudentSummary, 0)
			for _, id := range sess.StudentIDs {
				hexID := id.Hex()
				studentIDs = append(studentIDs, hexID)
				if sum, exists := studentSummaryMap[hexID]; exists {
					studentSummaries = append(studentSummaries, sum)
				}
			}

			occurrences = append(occurrences, SessionOccurrenceResponse{
				InstanceID:      fmt.Sprintf("inst-%s-%s", sess.ID.Hex(), dateStr),
				BlockID:         sess.ID.Hex(),
				Day:             dayOfWeek,
				Date:            dateStr,
				StartTime:       sess.StartTime,
				EndTime:         sess.EndTime,
				DurationMinutes: sess.DurationMinutes,
				Type:            sess.Type,
				StudentIDs:      studentIDs,
				Students:        studentSummaries,
				IsRecurring:     sess.IsRecurring,
				IsException:     false,
				Notes:           sess.Notes,
			})
		}

		// 2. Add sessions that were rescheduled TO this date
		if movedExceptions, ok := rescheduledMovedToDate[dateStr]; ok {
			for _, exc := range movedExceptions {
				// Find parent session
				var parentSession *Session
				for i := range sessions {
					if sessions[i].ID == exc.SessionID {
						parentSession = &sessions[i]
						break
					}
				}
				if parentSession == nil {
					// If not in preloaded list (e.g. filtered by student), query DB
					s, err := s.repo.GetSessionByID(ctx, exc.SessionID)
					if err == nil && s.Status == SessionStatusActive {
						parentSession = s
					}
				}
				if parentSession == nil || parentSession.Status != SessionStatusActive {
					continue
				}

				if studentOID != nil {
					hasStudent := false
					for _, sid := range parentSession.StudentIDs {
						if sid == *studentOID {
							hasStudent = true
							break
						}
					}
					if !hasStudent {
						continue
					}
				}

				startTime := parentSession.StartTime
				if exc.NewStartTime != nil && *exc.NewStartTime != "" {
					startTime = *exc.NewStartTime
				}
				endTime := parentSession.EndTime
				if exc.NewEndTime != nil && *exc.NewEndTime != "" {
					endTime = *exc.NewEndTime
				}

				sStart, _ := TimeToMinutes(startTime)
				sEnd, _ := TimeToMinutes(endTime)
				duration := sEnd - sStart

				studentIDs := make([]string, 0, len(parentSession.StudentIDs))
				studentSummaries := make([]StudentSummary, 0)
				for _, id := range parentSession.StudentIDs {
					hexID := id.Hex()
					studentIDs = append(studentIDs, hexID)
					if sum, exists := studentSummaryMap[hexID]; exists {
						studentSummaries = append(studentSummaries, sum)
					}
				}

				occurrences = append(occurrences, SessionOccurrenceResponse{
					InstanceID:      fmt.Sprintf("inst-exc-%s", exc.ID.Hex()),
					BlockID:         parentSession.ID.Hex(),
					Day:             dayOfWeek,
					Date:            dateStr,
					StartTime:       startTime,
					EndTime:         endTime,
					DurationMinutes: duration,
					Type:            parentSession.Type,
					StudentIDs:      studentIDs,
					Students:        studentSummaries,
					IsRecurring:     true,
					IsException:     true,
					ExceptionType:   ExceptionTypeRescheduled,
					ExceptionID:     exc.ID.Hex(),
					Notes:           exc.Notes,
				})
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Sort chronologically by date then by startTime
	sort.Slice(occurrences, func(i, j int) bool {
		if occurrences[i].Date != occurrences[j].Date {
			return occurrences[i].Date < occurrences[j].Date
		}
		startI, _ := TimeToMinutes(occurrences[i].StartTime)
		startJ, _ := TimeToMinutes(occurrences[j].StartTime)
		return startI < startJ
	})

	return occurrences, nil
}

// ==========================================
// CALENDAR EXCEPTIONS METHODS
// ==========================================

func (s *service) CreateException(ctx context.Context, sessionIDStr string, req CreateExceptionRequest) (*ExceptionResponse, error) {
	sessionOID, err := ParseHexID(sessionIDStr)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	session, err := s.repo.GetSessionByID(ctx, sessionOID)
	if err != nil {
		return nil, err
	}

	date := strings.TrimSpace(req.Date)
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, ErrInvalidDateFormat
	}

	// Check if an exception already exists for this session and date
	existing, err := s.repo.GetExceptionBySessionAndDate(ctx, sessionOID, date)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateException
	}

	excType := strings.ToLower(strings.TrimSpace(req.Type))
	if excType != ExceptionTypeCancelled && excType != ExceptionTypeRescheduled {
		return nil, ErrInvalidExceptionType
	}

	var newDate *string
	var newDay *DayOfWeek
	var newStart *string
	var newEnd *string

	if excType == ExceptionTypeRescheduled {
		if req.NewDate == nil || req.NewStartTime == nil || req.NewEndTime == nil {
			return nil, ErrRescheduledMissingData
		}
		nd := strings.TrimSpace(*req.NewDate)
		parsedNewDate, err := time.ParseInLocation("2006-01-02", nd, s.loc)
		if err != nil {
			return nil, fmt.Errorf("%w: nueva fecha debe ser YYYY-MM-DD", ErrInvalidDateFormat)
		}
		newDate = &nd

		d := TimeWeekdayToDayOfWeek(parsedNewDate.Weekday())
		newDay = &d

		ns := strings.TrimSpace(*req.NewStartTime)
		ne := strings.TrimSpace(*req.NewEndTime)
		startM, err := TimeToMinutes(ns)
		if err != nil {
			return nil, fmt.Errorf("%w: hora de inicio '%s'", ErrInvalidTimeFormat, ns)
		}
		endM, err := TimeToMinutes(ne)
		if err != nil {
			return nil, fmt.Errorf("%w: hora de término '%s'", ErrInvalidTimeFormat, ne)
		}
		if endM <= startM {
			return nil, ErrEndTimeBeforeStartTime
		}
		if endM-startM < 15 {
			return nil, ErrSessionDurationTooShort
		}

		// Operating hours
		opStart, _ := TimeToMinutes(OperatingStart)
		opEnd, _ := TimeToMinutes(OperatingEnd)
		if startM < opStart || endM > opEnd {
			return nil, ErrOutsideOperatingHours
		}

		// Validate availability of destination day
		destAvail, err := s.repo.GetAvailabilityByDay(ctx, *newDay)
		if err != nil {
			return nil, err
		}
		if destAvail == nil || !destAvail.Enabled || len(destAvail.Intervals) == 0 {
			return nil, fmt.Errorf("%w para el día de destino %s", ErrSessionOutsideAvailability, *newDay)
		}
		fits := false
		for _, interval := range destAvail.Intervals {
			iStart, _ := TimeToMinutes(interval.StartTime)
			iEnd, _ := TimeToMinutes(interval.EndTime)
			if FitsInInterval(startM, endM, iStart, iEnd) {
				fits = true
				break
			}
		}
		if !fits {
			return nil, fmt.Errorf("%w (%s — %s queda fuera de la disponibilidad para el %s)", ErrSessionOutsideAvailability, ns, ne, *newDay)
		}

		// Check conflict with destination sessions
		destOccurrences, err := s.GetOccurrences(ctx, SessionsFilter{From: nd, To: nd})
		if err != nil {
			return nil, err
		}
		for _, occ := range destOccurrences {
			if occ.BlockID == session.ID.Hex() && occ.Date == date {
				continue // ignore the session occurrence being rescheduled
			}
			occStart, _ := TimeToMinutes(occ.StartTime)
			occEnd, _ := TimeToMinutes(occ.EndTime)
			if IntervalsOverlap(startM, endM, occStart, occEnd) {
				return nil, fmt.Errorf("%w con la sesión existente en la fecha de destino (%s — %s)", ErrSessionOverlap, occ.StartTime, occ.EndTime)
			}
		}

		newStart = &ns
		newEnd = &ne
	}

	notes := strings.TrimSpace(req.Notes)
	if notes == "" {
		if excType == ExceptionTypeCancelled {
			notes = fmt.Sprintf("Cancelación puntual para la fecha %s.", date)
		} else {
			notes = fmt.Sprintf("Reagendado para el %s %s a %s.", *newDate, *newStart, *newEnd)
		}
	}

	exception := &ScheduleException{
		ID:           primitive.NewObjectID(),
		SessionID:    sessionOID,
		OriginalDate: date,
		Type:         excType,
		NewDate:      newDate,
		NewDay:       newDay,
		NewStartTime: newStart,
		NewEndTime:   newEnd,
		Notes:        notes,
	}

	if err := s.repo.CreateException(ctx, exception); err != nil {
		return nil, err
	}

	resp := ToExceptionResponse(exception)
	return &resp, nil
}

func (s *service) ListExceptions(ctx context.Context, filter ExceptionsFilter) ([]ExceptionResponse, error) {
	from := strings.TrimSpace(filter.From)
	to := strings.TrimSpace(filter.To)

	if from == "" {
		from = "1970-01-01"
	}
	if to == "" {
		to = "2099-12-31"
	}

	list, err := s.repo.ListExceptionsByPeriod(ctx, from, to)
	if err != nil {
		return nil, err
	}

	res := make([]ExceptionResponse, 0, len(list))
	parentSessionActive := make(map[primitive.ObjectID]bool)

	for _, exc := range list {
		active, checked := parentSessionActive[exc.SessionID]
		if !checked {
			parent, err := s.repo.GetSessionByID(ctx, exc.SessionID)
			active = (err == nil && parent != nil && parent.Status == SessionStatusActive)
			parentSessionActive[exc.SessionID] = active
		}
		if !active {
			continue
		}
		res = append(res, ToExceptionResponse(&exc))
	}
	return res, nil
}

func (s *service) DeleteException(ctx context.Context, sessionIDStr, exceptionIDStr string) error {
	excOID, err := ParseHexID(exceptionIDStr)
	if err != nil {
		return ErrExceptionNotFound
	}
	return s.repo.DeleteException(ctx, excOID)
}

// ==========================================
// INTERNAL HELPERS
// ==========================================

func (s *service) validateStudents(ctx context.Context, studentIDs []string) ([]primitive.ObjectID, []StudentSummary, error) {
	if len(studentIDs) == 0 {
		return []primitive.ObjectID{}, []StudentSummary{}, nil
	}

	// Check duplicates
	seen := make(map[string]struct{})
	var oids []primitive.ObjectID
	for _, id := range studentIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			return nil, nil, ErrDuplicateStudentInSession
		}
		seen[trimmed] = struct{}{}

		oid, err := ParseHexID(trimmed)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: identificador de alumno inválido '%s'", ErrStudentNotFoundOrArchived, trimmed)
		}
		oids = append(oids, oid)
	}

	var summaries []StudentSummary
	if s.studentReader != nil && len(studentIDs) > 0 {
		sums, err := s.studentReader.GetStudentsSummary(ctx, studentIDs)
		if err != nil {
			return nil, nil, err
		}
		if len(sums) != len(oids) {
			return nil, nil, ErrStudentNotFoundOrArchived
		}
		for _, sum := range sums {
			if sum.Archived {
				return nil, nil, fmt.Errorf("%w: alumno '%s' archivado", ErrStudentNotFoundOrArchived, sum.FullName)
			}
		}
		summaries = sums
	}

	return oids, summaries, nil
}

// checkSessionConflicts checks whether a session [startM, endM] on day with date range [startDate, endDate]
// conflicts with any active recurring session or calendar exception.
func (s *service) checkSessionConflicts(
	ctx context.Context,
	excludeID primitive.ObjectID,
	day DayOfWeek,
	startM, endM int,
	startDate string,
	endDate *string,
) error {
	// 1. Check overlap with other active sessions on the same day
	activeSessions, err := s.repo.ListActiveSessions(ctx, &day, nil)
	if err != nil {
		return err
	}

	for _, other := range activeSessions {
		if !excludeID.IsZero() && other.ID == excludeID {
			continue
		}

		otherStartM, err1 := TimeToMinutes(other.StartTime)
		otherEndM, err2 := TimeToMinutes(other.EndTime)
		if err1 != nil || err2 != nil {
			continue
		}

		if !IntervalsOverlap(startM, endM, otherStartM, otherEndM) {
			continue
		}

		// Time of day overlaps. Now check if date ranges overlap on at least one real occurrence of day.
		coincide, firstConflictDate := s.sessionsCoincideOnDate(day, startDate, endDate, other.StartDate, other.EndDate)
		if !coincide {
			continue
		}

		// If both are infinite into the future, collision is inevitable
		if endDate == nil && other.EndDate == nil {
			return fmt.Errorf("%w con la sesión existente (%s — %s)", ErrSessionOverlap, other.StartTime, other.EndTime)
		}

		// For finite overlaps, check if all collisions are excepted
		hasCollision, conflictDate := s.hasActiveCollision(ctx, other.ID, day, firstConflictDate, endDate, other.EndDate)
		if hasCollision {
			return fmt.Errorf("%w con la sesión existente en la fecha %s (%s — %s)", ErrSessionOverlap, conflictDate, other.StartTime, other.EndTime)
		}
	}

	// 2. Also check if any other session was rescheduled into this day and time on a date within [startDate, endDate]
	toQuery := "2099-12-31"
	if endDate != nil && *endDate != "" {
		toQuery = *endDate
	}
	exceptions, err := s.repo.ListExceptionsByPeriod(ctx, startDate, toQuery)
	if err == nil {
		for _, exc := range exceptions {
			if !excludeID.IsZero() && exc.SessionID == excludeID {
				continue
			}
			if exc.Type != ExceptionTypeRescheduled || exc.NewDate == nil || exc.NewDay == nil {
				continue
			}
			if *exc.NewDay != day {
				continue
			}
			if *exc.NewDate < startDate {
				continue
			}
			if endDate != nil && *exc.NewDate > *endDate {
				continue
			}

			if exc.NewStartTime != nil && exc.NewEndTime != nil {
				excStartM, err1 := TimeToMinutes(*exc.NewStartTime)
				excEndM, err2 := TimeToMinutes(*exc.NewEndTime)
				if err1 == nil && err2 == nil {
					if IntervalsOverlap(startM, endM, excStartM, excEndM) {
						// Ignorar si la sesión principal fue archivada o no existe
						parent, err := s.repo.GetSessionByID(ctx, exc.SessionID)
						if err != nil || parent == nil || parent.Status != SessionStatusActive {
							continue
						}
						return fmt.Errorf("%w con una sesión reprogramada para el %s (%s — %s)", ErrSessionOverlap, *exc.NewDate, *exc.NewStartTime, *exc.NewEndTime)
					}
				}
			}
		}
	}

	return nil
}

func (s *service) sessionsCoincideOnDate(
	day DayOfWeek,
	startA string,
	endA *string,
	startB string,
	endB *string,
) (bool, string) {
	interStart := startA
	if startB > interStart {
		interStart = startB
	}

	var interEnd *string
	if endA != nil && endB != nil {
		if *endA < *endB {
			interEnd = endA
		} else {
			interEnd = endB
		}
	} else if endA != nil {
		interEnd = endA
	} else if endB != nil {
		interEnd = endB
	}

	if interEnd != nil && interStart > *interEnd {
		return false, ""
	}

	parsedStart, err := time.ParseInLocation("2006-01-02", interStart, s.loc)
	if err != nil {
		return false, ""
	}
	targetWeekday, err := DayOfWeekToTimeWeekday(day)
	if err != nil {
		return false, ""
	}

	daysUntil := (int(targetWeekday) - int(parsedStart.Weekday()) + 7) % 7
	firstOcc := parsedStart.AddDate(0, 0, daysUntil)
	firstOccStr := firstOcc.Format("2006-01-02")

	if interEnd != nil && firstOccStr > *interEnd {
		return false, ""
	}

	return true, firstOccStr
}

func (s *service) hasActiveCollision(
	ctx context.Context,
	otherSessionID primitive.ObjectID,
	day DayOfWeek,
	firstConflictDate string,
	endA *string,
	endB *string,
) (bool, string) {
	var maxDate *string
	if endA != nil && endB != nil {
		if *endA < *endB {
			maxDate = endA
		} else {
			maxDate = endB
		}
	} else if endA != nil {
		maxDate = endA
	} else if endB != nil {
		maxDate = endB
	}

	currDate, err := time.ParseInLocation("2006-01-02", firstConflictDate, s.loc)
	if err != nil {
		return true, firstConflictDate
	}

	for i := 0; i < 100; i++ {
		dateStr := currDate.Format("2006-01-02")
		if maxDate != nil && dateStr > *maxDate {
			break
		}

		exc, err := s.repo.GetExceptionBySessionAndDate(ctx, otherSessionID, dateStr)
		if err == nil && exc != nil {
			if exc.Type == ExceptionTypeCancelled || exc.Type == ExceptionTypeRescheduled {
				currDate = currDate.AddDate(0, 0, 7)
				continue
			}
		}

		return true, dateStr
	}

	return false, ""
}
