package students

import (
	"strings"
	"time"
)

// CreateStudentRequest representa el payload recibido en POST /api/v1/students.
type CreateStudentRequest struct {
	FullName            string             `json:"fullName"`
	Frequency           string             `json:"frequency"`
	Email               string             `json:"email"`
	Phone               string             `json:"phone"`
	PlanName            string             `json:"planName"`
	MembershipStatus    string             `json:"membershipStatus"`
	FixedSchedule       string             `json:"fixedSchedule"`
	MembershipStartDate string             `json:"membershipStartDate"`
	MembershipDueDate   string             `json:"membershipDueDate"`
	Notes               string             `json:"notes"`
	HealthConditions    string             `json:"healthConditions"`
	EmergencyContact    string             `json:"emergencyContact"`
	Onboarding          *StudentOnboarding `json:"onboarding,omitempty"`
}

// UpdateStudentRequest representa el payload recibido en PUT /api/v1/students/:id.
// Se utilizan punteros para permitir modificaciones parciales seguras.
type UpdateStudentRequest struct {
	FullName            *string            `json:"fullName"`
	Frequency           *string            `json:"frequency"`
	Email               *string            `json:"email"`
	Phone               *string            `json:"phone"`
	PlanName            *string            `json:"planName"`
	MembershipStatus    *string            `json:"membershipStatus"`
	FixedSchedule       *string            `json:"fixedSchedule"`
	MembershipStartDate *string            `json:"membershipStartDate"`
	MembershipDueDate   *string            `json:"membershipDueDate"`
	Notes               *string            `json:"notes"`
	HealthConditions    *string            `json:"healthConditions"`
	EmergencyContact    *string            `json:"emergencyContact"`
	Onboarding          *StudentOnboarding `json:"onboarding,omitempty"`
}

// StudentResponse representa el objeto JSON que se expone al frontend Vue.
// Es 100% compatible con la interfaz TypeScript 'Student'.
type StudentResponse struct {
	ID                  string             `json:"id"`
	FullName            string             `json:"fullName"`
	Email               string             `json:"email"`
	Phone               string             `json:"phone"`
	PlanName            string             `json:"planName"`
	Frequency           string             `json:"frequency"`
	MembershipStatus    string             `json:"membershipStatus"`
	FixedSchedule       string             `json:"fixedSchedule"`
	MembershipStartDate string             `json:"membershipStartDate"`
	MembershipDueDate   string             `json:"membershipDueDate"`
	Notes               string             `json:"notes"`
	HealthConditions    string             `json:"healthConditions"`
	EmergencyContact    string             `json:"emergencyContact"`
	Onboarding          *StudentOnboarding `json:"onboarding,omitempty"`
	CreatedAt           string             `json:"createdAt"`
	UpdatedAt           string             `json:"updatedAt"`
}

// StudentListFilter contiene los criterios de búsqueda opcionales para GET /api/v1/students.
type StudentListFilter struct {
	Name      string
	Status    string
	Frequency string
}

// ToStudentResponse transforma un modelo de persistencia Student en su correspondiente DTO para respuesta HTTP.
func ToStudentResponse(s *Student) StudentResponse {
	createdAtStr := s.CreatedAt.Format(time.RFC3339)
	if s.CreatedAt.IsZero() {
		createdAtStr = time.Now().UTC().Format(time.RFC3339)
	}

	updatedAtStr := s.UpdatedAt.Format(time.RFC3339)
	if s.UpdatedAt.IsZero() {
		updatedAtStr = createdAtStr
	}

	return StudentResponse{
		ID:                  s.ID.Hex(),
		FullName:            strings.TrimSpace(s.FullName),
		Email:               strings.TrimSpace(s.Email),
		Phone:               strings.TrimSpace(s.Phone),
		PlanName:            strings.TrimSpace(s.PlanName),
		Frequency:           strings.TrimSpace(s.Frequency),
		MembershipStatus:    strings.TrimSpace(s.MembershipStatus),
		FixedSchedule:       strings.TrimSpace(s.FixedSchedule),
		MembershipStartDate: strings.TrimSpace(s.MembershipStartDate),
		MembershipDueDate:   strings.TrimSpace(s.MembershipDueDate),
		Notes:               strings.TrimSpace(s.Notes),
		HealthConditions:    strings.TrimSpace(s.HealthConditions),
		EmergencyContact:    strings.TrimSpace(s.EmergencyContact),
		Onboarding:          s.Onboarding,
		CreatedAt:           createdAtStr,
		UpdatedAt:           updatedAtStr,
	}
}

// ToStudentResponseList convierte un slice de entidades Student en un slice de DTOs.
func ToStudentResponseList(list []Student) []StudentResponse {
	res := make([]StudentResponse, len(list))
	for i := range list {
		res[i] = ToStudentResponse(&list[i])
	}
	return res
}
