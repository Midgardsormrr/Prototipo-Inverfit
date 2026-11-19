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

	StatusActive   = "Activo"
	StatusToReview = "Por revisar"
	StatusInactive = "Inactivo"
)

// StudentOnboarding almacena las 20 preguntas oficiales del formulario de Pablo + consentimiento y firma.
type StudentOnboarding struct {
	// 1. Demográficos
	Rut             string `bson:"rut,omitempty" json:"rut"`
	BirthDate       string `bson:"birthDate,omitempty" json:"birthDate"`
	Age             int    `bson:"age,omitempty" json:"age"`
	Gender          string `bson:"gender,omitempty" json:"gender"`
	Comuna          string `bson:"comuna,omitempty" json:"comuna"`
	Occupation      string `bson:"occupation,omitempty" json:"occupation"`
	HasChildren     string `bson:"hasChildren,omitempty" json:"hasChildren"`
	ChildrenCount   string `bson:"childrenCount,omitempty" json:"childrenCount"`

	// 2. Perfil Deportivo y Objetivos
	PrimaryGoal         string `bson:"primaryGoal,omitempty" json:"primaryGoal"`
	GoalDetails         string `bson:"goalDetails,omitempty" json:"goalDetails"`
	HasTrainedBefore    string `bson:"hasTrainedBefore,omitempty" json:"hasTrainedBefore"`
	TrainingExperience  string `bson:"trainingExperience,omitempty" json:"trainingExperience"`
	SportsPracticed     string `bson:"sportsPracticed,omitempty" json:"sportsPracticed"`
	AvailableDaysWeekly string `bson:"availableDaysWeekly,omitempty" json:"availableDaysWeekly"`

	// 3. Ficha Clínica y Médica (PAR-Q+)
	HadCovid                string `bson:"hadCovid,omitempty" json:"hadCovid"`
	CovidSequelae           string `bson:"covidSequelae,omitempty" json:"covidSequelae"`
	ChronicDiseases         string `bson:"chronicDiseases,omitempty" json:"chronicDiseases"`
	SurgeriesInjuries       string `bson:"surgeriesInjuries,omitempty" json:"surgeriesInjuries"`
	InjuriesPreventExercise string `bson:"injuriesPreventExercise,omitempty" json:"injuriesPreventExercise"`
	CurrentDifficulties     string `bson:"currentDifficulties,omitempty" json:"currentDifficulties"`
	FrequentConsumption     string `bson:"frequentConsumption,omitempty" json:"frequentConsumption"`

	// 4. Consentimiento y Firma Digital
	ConsentAccepted  bool   `bson:"consentAccepted" json:"consentAccepted"`
	ConsentTimestamp string `bson:"consentTimestamp,omitempty" json:"consentTimestamp"`
	DigitalSignature string `bson:"digitalSignature,omitempty" json:"digitalSignature"`
	TermsVersion     string `bson:"termsVersion,omitempty" json:"termsVersion"`
}

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
	Onboarding          *StudentOnboarding `bson:"onboarding,omitempty"`
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
