package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DayOfWeek defines the day of the week in canonical Spanish short format ('Lun', 'Mar', etc.).
type DayOfWeek string

const (
	DayLun DayOfWeek = "Lun"
	DayMar DayOfWeek = "Mar"
	DayMie DayOfWeek = "Mié"
	DayJue DayOfWeek = "Jue"
	DayVie DayOfWeek = "Vie"
	DaySab DayOfWeek = "Sáb"
	DayDom DayOfWeek = "Dom"
)

// CanonicalDays lists all 7 days of the week in standard Monday-first order.
var CanonicalDays = []DayOfWeek{DayLun, DayMar, DayMie, DayJue, DayVie, DaySab, DayDom}

// Operating hours limits for INVERFIT gym
const (
	OperatingStart = "06:30"
	OperatingEnd   = "22:00"
	DefaultTimezone = "America/Santiago"
)

// Session statuses
const (
	SessionStatusActive   = "active"
	SessionStatusArchived = "archived"
)

// Exception types
const (
	ExceptionTypeCancelled   = "cancelled"
	ExceptionTypeRescheduled = "rescheduled"
)

// Session types
const (
	SessionTypeIndividual = "Individual"
	SessionTypeGroup      = "Sesión grupal"
)

// NormalizeDayOfWeek accepts diverse representations of weekdays and returns canonical DayOfWeek.
func NormalizeDayOfWeek(input string) (DayOfWeek, error) {
	cleaned := strings.ToLower(strings.TrimSpace(input))
	cleaned = strings.ReplaceAll(cleaned, "á", "a")
	cleaned = strings.ReplaceAll(cleaned, "é", "e")
	cleaned = strings.ReplaceAll(cleaned, "í", "i")
	cleaned = strings.ReplaceAll(cleaned, "ó", "o")
	cleaned = strings.ReplaceAll(cleaned, "ú", "u")

	switch cleaned {
	case "lun", "lunes", "mon", "monday":
		return DayLun, nil
	case "mar", "martes", "tue", "tuesday":
		return DayMar, nil
	case "mie", "miercoles", "wed", "wednesday":
		return DayMie, nil
	case "jue", "jueves", "thu", "thursday":
		return DayJue, nil
	case "vie", "viernes", "fri", "friday":
		return DayVie, nil
	case "sab", "sabado", "sat", "saturday":
		return DaySab, nil
	case "dom", "domingo", "sun", "sunday":
		return DayDom, nil
	default:
		return "", ErrInvalidDayOfWeek
	}
}

// TimeWeekdayToDayOfWeek maps Go's time.Weekday (0=Sun, 1=Mon, ..., 6=Sat) to canonical DayOfWeek.
func TimeWeekdayToDayOfWeek(w time.Weekday) DayOfWeek {
	switch w {
	case time.Monday:
		return DayLun
	case time.Tuesday:
		return DayMar
	case time.Wednesday:
		return DayMie
	case time.Thursday:
		return DayJue
	case time.Friday:
		return DayVie
	case time.Saturday:
		return DaySab
	case time.Sunday:
		return DayDom
	default:
		return DayLun
	}
}

// DayOfWeekToTimeWeekday maps canonical DayOfWeek to Go's time.Weekday.
func DayOfWeekToTimeWeekday(d DayOfWeek) (time.Weekday, error) {
	switch d {
	case DayLun:
		return time.Monday, nil
	case DayMar:
		return time.Tuesday, nil
	case DayMie:
		return time.Wednesday, nil
	case DayJue:
		return time.Thursday, nil
	case DayVie:
		return time.Friday, nil
	case DaySab:
		return time.Saturday, nil
	case DayDom:
		return time.Sunday, nil
	default:
		return time.Monday, ErrInvalidDayOfWeek
	}
}

// TimeToMinutes converts "HH:mm" to total minutes from midnight (0 to 1439).
func TimeToMinutes(timeStr string) (int, error) {
	parts := strings.Split(strings.TrimSpace(timeStr), ":")
	if len(parts) != 2 {
		return 0, ErrInvalidTimeFormat
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, ErrInvalidTimeFormat
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, ErrInvalidTimeFormat
	}
	return h*60 + m, nil
}

// MinutesToTime converts minutes from midnight to "HH:mm" string.
func MinutesToTime(minutes int) string {
	if minutes < 0 {
		minutes = 0
	}
	if minutes >= 24*60 {
		minutes = 24*60 - 1
	}
	h := minutes / 60
	m := minutes % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

// IntervalsOverlap checks if two time intervals [startA, endA] and [startB, endB] overlap.
// Rule: startA < endB && startB < endA
// Note: If startB == endA (consecutive), it evaluates to false (no overlap, valid!).
func IntervalsOverlap(startA, endA, startB, endB int) bool {
	return startA < endB && startB < endA
}

// FitsInInterval checks if session [sessionStart, sessionEnd] is completely inside interval [intStart, intEnd].
func FitsInInterval(sessionStart, sessionEnd, intStart, intEnd int) bool {
	return sessionStart >= intStart && sessionEnd <= intEnd
}

// AvailabilityInterval represents a single working window within a day.
type AvailabilityInterval struct {
	ID        string `bson:"id" json:"id"`
	StartTime string `bson:"start_time" json:"startTime"` // "HH:mm"
	EndTime   string `bson:"end_time" json:"endTime"`     // "HH:mm"
}

// DayAvailability represents the coach's scheduled working intervals for a specific day of the week.
type DayAvailability struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Day       DayOfWeek              `bson:"day" json:"day"`
	Enabled   bool                   `bson:"enabled" json:"enabled"`
	Intervals []AvailabilityInterval `bson:"intervals" json:"intervals"`
	Timezone  string                 `bson:"timezone" json:"timezone"`
	CreatedAt time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updatedAt"`
}

// Session represents a recurring schedule block assigned to one or more students.
type Session struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Day             DayOfWeek            `bson:"day" json:"day"`
	StartTime       string               `bson:"start_time" json:"startTime"`       // "HH:mm"
	EndTime         string               `bson:"end_time" json:"endTime"`           // "HH:mm"
	DurationMinutes int                  `bson:"duration_minutes" json:"durationMinutes"`
	Type            string               `bson:"type" json:"type"`                 // "Individual" | "Sesión grupal"
	StudentIDs      []primitive.ObjectID `bson:"student_ids" json:"studentIds"`
	IsRecurring     bool                 `bson:"is_recurring" json:"isRecurring"`
	StartDate       string               `bson:"start_date" json:"startDate"`     // "YYYY-MM-DD"
	EndDate         *string              `bson:"end_date,omitempty" json:"endDate,omitempty"`
	Status          string               `bson:"status" json:"status"`             // "active" | "archived"
	Notes           string               `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
	ArchivedAt      *time.Time           `bson:"archived_at,omitempty" json:"archivedAt,omitempty"`
}

// ScheduleException represents a punctual modification (cancelled or rescheduled) of a recurring session occurrence.
type ScheduleException struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID    primitive.ObjectID `bson:"session_id" json:"sessionId"`
	OriginalDate string             `bson:"original_date" json:"originalDate"` // "YYYY-MM-DD"
	Type         string             `bson:"type" json:"type"`                   // "cancelled" | "rescheduled"
	NewDate      *string            `bson:"new_date,omitempty" json:"newDate,omitempty"` // "YYYY-MM-DD"
	NewDay       *DayOfWeek         `bson:"new_day,omitempty" json:"newDay,omitempty"`   // "Lun"
	NewStartTime *string            `bson:"new_start_time,omitempty" json:"newStartTime,omitempty"` // "HH:mm"
	NewEndTime   *string            `bson:"new_end_time,omitempty" json:"newEndTime,omitempty"`     // "HH:mm"
	Notes        string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}
