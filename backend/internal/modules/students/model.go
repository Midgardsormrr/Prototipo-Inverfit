package students

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constantes de valores permitidos para frecuencia y estado
const (
	Frequency1x = "1 vez/sem."
	Frequency2x = "2 veces/sem."
	Frequency3x = "3 veces/sem."
	Frequency4x = "4 veces/sem."

	StatusActive    = "Activo"
	StatusToReview  = "Por revisar"
	StatusInactive  = "Inactivo"
)

// Student representa el documento de persistencia en la colección MongoDB 'students'.
type Student struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	FullName            string             `bson:"fullName"`
	Email               string             `bson:"email"`
	Phone               string             `bson:"phone"`
	PlanName            string             `bson:"planName"`
	Frequency           string             `bson:"frequency"`
	MembershipStatus    string             `bson:"membershipStatus"`
	FixedSchedule       string             `bson:"fixedSchedule"`
	MembershipStartDate string             `bson:"membershipStartDate"`
	MembershipDueDate   string             `bson:"membershipDueDate"`
	Notes               string             `bson:"notes"`
	HealthConditions    string             `bson:"healthConditions"`
	EmergencyContact    string             `bson:"emergencyContact"`
	CreatedAt           time.Time          `bson:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt"`
	ArchivedAt          *time.Time         `bson:"archivedAt,omitempty"`
}

// IsValidFrequency valida que la frecuencia coincida con los valores soportados.
func IsValidFrequency(f string) bool {
	switch f {
	case Frequency1x, Frequency2x, Frequency3x, Frequency4x:
		return true
	default:
		return false
	}
}

// IsValidStatus valida que el estado de membresía coincida con los valores soportados.
func IsValidStatus(s string) bool {
	switch s {
	case StatusActive, StatusToReview, StatusInactive:
		return true
	default:
		return false
	}
}
