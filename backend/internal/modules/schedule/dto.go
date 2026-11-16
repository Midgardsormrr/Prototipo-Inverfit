package schedule

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateAvailabilityRequest holds the payload to set availability intervals for a day.
type UpdateAvailabilityRequest struct {
	Enabled   *bool                  `json:"enabled"`
	Intervals []AvailabilityInterval `json:"intervals"`
}

// DayAvailabilityResponse represents availability data returned by the API.
type DayAvailabilityResponse struct {
	Day       DayOfWeek              `json:"day"`
	Enabled   bool                   `json:"enabled"`
	Intervals []AvailabilityInterval `json:"intervals"`
	Timezone  string                 `json:"timezone"`
}

// AvailabilityConflictResponse is returned when changing availability conflicts with existing active sessions.
type AvailabilityConflictResponse struct {
	Message   string            `json:"message"`
	Conflicts []SessionResponse `json:"conflicts"`
}

// CreateSessionRequest defines the input payload for creating a recurring session.
type CreateSessionRequest struct {
	Day             string   `json:"day" binding:"required"`
	StartTime       string   `json:"startTime" binding:"required"`
	EndTime         string   `json:"endTime" binding:"required"`
	DurationMinutes *int     `json:"durationMinutes,omitempty"`
	StudentIDs      []string `json:"studentIds"`
	IsRecurring     *bool    `json:"isRecurring,omitempty"`
	StartDate       string   `json:"startDate,omitempty"`
	EndDate         *string  `json:"endDate,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}

// UpdateSessionRequest defines the payload for updating an existing recurring session.
type UpdateSessionRequest struct {
	Day         *string   `json:"day,omitempty"`
	StartTime   *string   `json:"startTime,omitempty"`
	EndTime     *string   `json:"endTime,omitempty"`
	StudentIDs  *[]string `json:"studentIds,omitempty"`
	IsRecurring *bool     `json:"isRecurring,omitempty"`
	StartDate   *string   `json:"startDate,omitempty"`
	EndDate     *string   `json:"endDate,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

// SessionResponse represents the full recurring session definition.
type SessionResponse struct {
	ID              string           `json:"id"`
	Day             DayOfWeek        `json:"day"`
	StartTime       string           `json:"startTime"`
	EndTime         string           `json:"endTime"`
	DurationMinutes int              `json:"durationMinutes"`
	Type            string           `json:"type"`
	StudentIDs      []string         `json:"studentIds"`
	Students        []StudentSummary `json:"students,omitempty"`
	IsRecurring     bool             `json:"isRecurring"`
	StartDate       string           `json:"startDate"`
	EndDate         *string          `json:"endDate,omitempty"`
	Status          string           `json:"status"`
	Notes           string           `json:"notes"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

// SessionOccurrenceResponse represents a concrete session instance for a calendar date,
// projected from recurring rules and punctual calendar exceptions.
type SessionOccurrenceResponse struct {
	InstanceID      string           `json:"instanceId"`
	BlockID         string           `json:"blockId"` // Matches frontend blockId
	Day             DayOfWeek        `json:"day"`
	Date            string           `json:"date"` // "YYYY-MM-DD"
	StartTime       string           `json:"startTime"`
	EndTime         string           `json:"endTime"`
	DurationMinutes int              `json:"durationMinutes"`
	Type            string           `json:"type"`
	StudentIDs      []string         `json:"studentIds"`
	Students        []StudentSummary `json:"students,omitempty"`
	IsRecurring     bool             `json:"isRecurring"`
	IsException     bool             `json:"isException"`
	ExceptionType   string           `json:"exceptionType,omitempty"` // "cancelled" | "rescheduled"
	ExceptionID     string           `json:"exceptionId,omitempty"`
	Notes           string           `json:"notes"`
}

// CreateExceptionRequest defines input for creating a calendar exception.
type CreateExceptionRequest struct {
	Date         string  `json:"date" binding:"required"` // "YYYY-MM-DD"
	Type         string  `json:"type" binding:"required"` // "cancelled" | "rescheduled"
	NewDate      *string `json:"newDate,omitempty"`       // "YYYY-MM-DD"
	NewDay       *string `json:"newDay,omitempty"`        // "Lun", etc.
	NewStartTime *string `json:"newStartTime,omitempty"`  // "HH:mm"
	NewEndTime   *string `json:"newEndTime,omitempty"`    // "HH:mm"
	Notes        string  `json:"notes,omitempty"`
}

// ExceptionResponse represents a calendar exception.
type ExceptionResponse struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"sessionId"`
	OriginalDate string    `json:"originalDate"`
	Type         string    `json:"type"`
	NewDate      *string   `json:"newDate,omitempty"`
	NewDay       *string   `json:"newDay,omitempty"`
	NewStartTime *string   `json:"newStartTime,omitempty"`
	NewEndTime   *string   `json:"newEndTime,omitempty"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"createdAt"`
}

// SessionsFilter defines query parameters for listing sessions / occurrences.
type SessionsFilter struct {
	From      string // "YYYY-MM-DD"
	To        string // "YYYY-MM-DD"
	StudentID string
}

// ExceptionsFilter defines query parameters for listing exceptions.
type ExceptionsFilter struct {
	From string // "YYYY-MM-DD"
	To   string // "YYYY-MM-DD"
}

// ToSessionResponse converts a domain Session model to SessionResponse.
func ToSessionResponse(s *Session, students []StudentSummary) SessionResponse {
	studentIDs := make([]string, 0, len(s.StudentIDs))
	for _, id := range s.StudentIDs {
		studentIDs = append(studentIDs, id.Hex())
	}
	if students == nil {
		students = make([]StudentSummary, 0)
	}

	return SessionResponse{
		ID:              s.ID.Hex(),
		Day:             s.Day,
		StartTime:       s.StartTime,
		EndTime:         s.EndTime,
		DurationMinutes: s.DurationMinutes,
		Type:            s.Type,
		StudentIDs:      studentIDs,
		Students:        students,
		IsRecurring:     s.IsRecurring,
		StartDate:       s.StartDate,
		EndDate:         s.EndDate,
		Status:          s.Status,
		Notes:           s.Notes,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

// ToExceptionResponse converts a domain ScheduleException model to ExceptionResponse.
func ToExceptionResponse(e *ScheduleException) ExceptionResponse {
	var newDayStr *string
	if e.NewDay != nil {
		s := string(*e.NewDay)
		newDayStr = &s
	}

	return ExceptionResponse{
		ID:           e.ID.Hex(),
		SessionID:    e.SessionID.Hex(),
		OriginalDate: e.OriginalDate,
		Type:         e.Type,
		NewDate:      e.NewDate,
		NewDay:       newDayStr,
		NewStartTime: e.NewStartTime,
		NewEndTime:   e.NewEndTime,
		Notes:        e.Notes,
		CreatedAt:    e.CreatedAt,
	}
}

// ParseHexID validates and converts a 24-character hexadecimal string into an ObjectID.
func ParseHexID(hexStr string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hexStr)
}
