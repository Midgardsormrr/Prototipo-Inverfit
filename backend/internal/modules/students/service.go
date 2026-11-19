package students

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service define el contrato de lógica de negocio para el módulo de alumnos.
type Service interface {
	Create(ctx context.Context, req CreateStudentRequest) (*StudentResponse, error)
	GetByID(ctx context.Context, idStr string) (*StudentResponse, error)
	List(ctx context.Context, filter StudentListFilter) ([]StudentResponse, error)
	Update(ctx context.Context, idStr string, req UpdateStudentRequest) (*StudentResponse, error)
	Delete(ctx context.Context, idStr string) error
}

type service struct {
	repo Repository
}

// NewService inicializa una nueva instancia de Service con su repositorio correspondiente.
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create valida las reglas de negocio, aplica defaults defensivos y registra un nuevo alumno.
func (s *service) Create(ctx context.Context, req CreateStudentRequest) (*StudentResponse, error) {
	name := strings.TrimSpace(req.FullName)
	if name == "" {
		return nil, ErrEmptyName
	}

	freq := strings.TrimSpace(req.Frequency)
	if freq == "" || !IsValidFrequency(freq) {
		return nil, ErrInvalidFrequency
	}

	status := strings.TrimSpace(req.MembershipStatus)
	if status == "" {
		status = StatusActive
	} else if !IsValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		// Generar correo sugerido si no se proporciona
		sanitizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		email = sanitizedName + "@ejemplo.com"
	}

	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		phone = "+56 9 Por definir"
	}

	planName := strings.TrimSpace(req.PlanName)
	if planName == "" {
		planName = "Membresía Mensual"
	}

	fixedSchedule := strings.TrimSpace(req.FixedSchedule)
	if fixedSchedule == "" {
		fixedSchedule = "Días y horas por confirmar"
	}

	startDate := strings.TrimSpace(req.MembershipStartDate)
	if startDate == "" {
		startDate = time.Now().UTC().Format("2006-01-02")
	}

	dueDate := strings.TrimSpace(req.MembershipDueDate)
	if dueDate == "" {
		dueDate = "Por definir"
	}

	health := strings.TrimSpace(req.HealthConditions)
	if health == "" {
		health = "Sin observaciones de salud registradas"
	}

	emergency := strings.TrimSpace(req.EmergencyContact)
	if emergency == "" {
		emergency = "Por registrar"
	}

	student := &Student{
		FullName:            name,
		Email:               email,
		Phone:               phone,
		PlanName:            planName,
		Frequency:           freq,
		MembershipStatus:    status,
		FixedSchedule:       fixedSchedule,
		MembershipStartDate: startDate,
		MembershipDueDate:   dueDate,
		Notes:               strings.TrimSpace(req.Notes),
		HealthConditions:    health,
		EmergencyContact:    emergency,
		Onboarding:          req.Onboarding,
	}

	if req.Onboarding != nil {
		if health == "Sin observaciones de salud registradas" || health == "" {
			var healthParts []string
			if req.Onboarding.ChronicDiseases != "" && strings.ToLower(req.Onboarding.ChronicDiseases) != "no" && strings.ToLower(req.Onboarding.ChronicDiseases) != "ninguna" {
				healthParts = append(healthParts, "Crónicas: "+req.Onboarding.ChronicDiseases)
			}
			if req.Onboarding.SurgeriesInjuries != "" && strings.ToLower(req.Onboarding.SurgeriesInjuries) != "no" && strings.ToLower(req.Onboarding.SurgeriesInjuries) != "ninguna" {
				healthParts = append(healthParts, "Lesiones/Cirugías: "+req.Onboarding.SurgeriesInjuries)
			}
			if len(healthParts) > 0 {
				student.HealthConditions = strings.Join(healthParts, " · ")
			}
		}
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return nil, err
	}

	res := ToStudentResponse(student)
	return &res, nil
}

// GetByID busca un alumno por su ID en formato hex string.
func (s *service) GetByID(ctx context.Context, idStr string) (*StudentResponse, error) {
	oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	student, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}

	res := ToStudentResponse(student)
	return &res, nil
}

// List recupera los alumnos activos según los criterios de filtro opcionales.
func (s *service) List(ctx context.Context, filter StudentListFilter) ([]StudentResponse, error) {
	filter.Name = strings.TrimSpace(filter.Name)
	filter.Status = strings.TrimSpace(filter.Status)
	filter.Frequency = strings.TrimSpace(filter.Frequency)

	if filter.Status != "" && !IsValidStatus(filter.Status) {
		return nil, ErrInvalidStatus
	}

	if filter.Frequency != "" && !IsValidFrequency(filter.Frequency) {
		return nil, ErrInvalidFrequency
	}

	students, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return ToStudentResponseList(students), nil
}

// Update modifica los campos de un alumno existente respetando modificaciones parciales.
func (s *service) Update(ctx context.Context, idStr string, req UpdateStudentRequest) (*StudentResponse, error) {
	oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	student, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}

	if req.FullName != nil {
		name := strings.TrimSpace(*req.FullName)
		if name == "" {
			return nil, ErrEmptyName
		}
		student.FullName = name
	}

	if req.Frequency != nil {
		freq := strings.TrimSpace(*req.Frequency)
		if !IsValidFrequency(freq) {
			return nil, ErrInvalidFrequency
		}
		student.Frequency = freq
	}

	if req.MembershipStatus != nil {
		status := strings.TrimSpace(*req.MembershipStatus)
		if !IsValidStatus(status) {
			return nil, ErrInvalidStatus
		}
		student.MembershipStatus = status
	}

	if req.Email != nil {
		student.Email = strings.TrimSpace(*req.Email)
	}

	if req.Phone != nil {
		student.Phone = strings.TrimSpace(*req.Phone)
	}

	if req.PlanName != nil {
		student.PlanName = strings.TrimSpace(*req.PlanName)
	}

	if req.FixedSchedule != nil {
		student.FixedSchedule = strings.TrimSpace(*req.FixedSchedule)
	}

	if req.MembershipStartDate != nil {
		student.MembershipStartDate = strings.TrimSpace(*req.MembershipStartDate)
	}

	if req.MembershipDueDate != nil {
		student.MembershipDueDate = strings.TrimSpace(*req.MembershipDueDate)
	}

	if req.Notes != nil {
		student.Notes = strings.TrimSpace(*req.Notes)
	}

	if req.HealthConditions != nil {
		student.HealthConditions = strings.TrimSpace(*req.HealthConditions)
	}

	if req.EmergencyContact != nil {
		student.EmergencyContact = strings.TrimSpace(*req.EmergencyContact)
	}

	if req.Onboarding != nil {
		student.Onboarding = req.Onboarding
	}

	if err := s.repo.Update(ctx, student); err != nil {
		return nil, err
	}

	res := ToStudentResponse(student)
	return &res, nil
}

// Delete ejecuta la eliminación lógica (archivado) del alumno.
func (s *service) Delete(ctx context.Context, idStr string) error {
	oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return ErrInvalidID
	}

	return s.repo.SoftDelete(ctx, oid)
}
