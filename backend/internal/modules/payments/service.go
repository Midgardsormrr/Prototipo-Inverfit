package payments

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	santiagoLocation *time.Location
)

func init() {
	var err error
	santiagoLocation, err = time.LoadLocation("America/Santiago")
	if err != nil {
		santiagoLocation = time.FixedZone("CLT", -3*3600)
	}
}

// Service defines business operations for the payments module.
type Service interface {
	RegisterPayment(ctx context.Context, req RegisterPaymentRequest) (*PaymentResponse, error)
	GetByID(ctx context.Context, idStr string) (*PaymentResponse, error)
	List(ctx context.Context, filter PaymentFilter) ([]PaymentResponse, error)
	GetByStudentID(ctx context.Context, studentIDStr string) ([]PaymentResponse, error)
	Update(ctx context.Context, idStr string, req UpdatePaymentRequest) (*PaymentResponse, error)
	Archive(ctx context.Context, idStr string) error
	SimulateReminder(ctx context.Context, idStr string) (*PaymentResponse, error)
	SimulateReceipt(ctx context.Context, idStr string) (*PaymentResponse, error)
	GetFinancialSummary(ctx context.Context) (*FinancialSummaryResponse, error)
}

type service struct {
	repo          Repository
	studentReader StudentReader
}

// NewService creates a new payments service instance with required dependencies.
func NewService(repo Repository, studentReader StudentReader) Service {
	return &service{
		repo:          repo,
		studentReader: studentReader,
	}
}

// RegisterPayment handles manual payment registration, matching existing pending records or creating new ones.
func (s *service) RegisterPayment(ctx context.Context, req RegisterPaymentRequest) (*PaymentResponse, error) {
	stOID, err := primitive.ObjectIDFromHex(strings.TrimSpace(req.StudentID))
	if err != nil {
		return nil, ErrInvalidID
	}

	// 1. Verify student exists via decoupled reader
	var student *StudentInfo
	if s.studentReader != nil {
		student, err = s.studentReader.GetStudent(ctx, req.StudentID)
		if err != nil || student == nil {
			return nil, ErrStudentNotFound
		}
	} else {
		student = &StudentInfo{
			ID:               req.StudentID,
			FullName:         "Alumno",
			PlanName:         "Membresía Mensual",
			Frequency:        "2 veces/sem.",
			MembershipStatus: "Activo",
		}
	}

	// 2. Validate amount
	if req.Amount != nil && *req.Amount < 0 {
		return nil, ErrInvalidAmount
	}

	// 3. Validate payment method
	method := strings.TrimSpace(req.PaymentMethod)
	if method == "" {
		method = MethodTransfer
	} else if !IsValidPaymentMethod(method) {
		return nil, ErrInvalidPaymentMethod
	}

	// 4. Validate or default payment date in America/Santiago
	nowSantiago := time.Now().In(santiagoLocation)
	todayStr := nowSantiago.Format("2006-01-02")
	payDateStr := strings.TrimSpace(req.PaymentDate)
	if payDateStr == "" {
		payDateStr = todayStr
	} else if !dateRegex.MatchString(payDateStr) {
		return nil, ErrInvalidDateFormat
	}

	dueDateStr := strings.TrimSpace(req.DueDate)
	if dueDateStr == "" {
		dueDateStr = "Por definir"
	}

	// 5. Generate transaction reference if not supplied
	txRef := strings.TrimSpace(req.TransactionReference)
	if txRef == "" {
		txRef = fmt.Sprintf("BOL-%s-%04d", nowSantiago.Format("200601"), rand.Intn(10000))
	}

	notes := strings.TrimSpace(req.Notes)

	// 6. Check for an existing pending payment for this student
	pendingPayment, err := s.repo.GetPendingByStudentID(ctx, stOID)
	if err != nil {
		return nil, err
	}

	var saved *Payment
	if pendingPayment != nil {
		// Update existing pending payment
		pendingPayment.Status = StatusPaid
		if req.Amount != nil {
			pendingPayment.Amount = req.Amount
		}
		pendingPayment.PaymentDate = &payDateStr
		pendingPayment.PaymentMethod = method
		pendingPayment.TransactionReference = txRef
		if notes != "" {
			pendingPayment.Notes = notes
		} else if pendingPayment.Notes == "" {
			pendingPayment.Notes = "Pago confirmado y registrado por el entrenador."
		}

		saved, err = s.repo.Update(ctx, pendingPayment)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new paid payment record
		newPayment := &Payment{
			StudentID:            stOID,
			StudentName:          student.FullName,
			PlanName:             student.PlanName,
			Frequency:            student.Frequency,
			DueDate:              dueDateStr,
			Amount:               req.Amount,
			Status:               StatusPaid,
			PaymentDate:          &payDateStr,
			PaymentMethod:        method,
			TransactionReference: txRef,
			ReminderSentCount:    0,
			LastReminderAt:       nil,
			ReceiptSent:          false,
			ReceiptSentAt:        nil,
			Notes:                notes,
		}
		if newPayment.Notes == "" {
			newPayment.Notes = "Nuevo pago registrado manualmente."
		}

		saved, err = s.repo.Create(ctx, newPayment)
		if err != nil {
			return nil, err
		}
	}

	// 7. If student's membership status is "Por revisar", activate it
	if s.studentReader != nil && student.MembershipStatus == "Por revisar" {
		_ = s.studentReader.ActivateMembershipIfReviewing(ctx, req.StudentID)
	}

	resp := ToPaymentResponse(saved)
	return &resp, nil
}

// GetByID returns the payment matching idStr.
func (s *service) GetByID(ctx context.Context, idStr string) (*PaymentResponse, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := ToPaymentResponse(p)
	return &resp, nil
}

// List returns active payments filtered by criteria.
func (s *service) List(ctx context.Context, filter PaymentFilter) ([]PaymentResponse, error) {
	if filter.Status != "" && !IsValidStatus(filter.Status) {
		return nil, ErrInvalidStatus
	}

	payments, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		result[i] = ToPaymentResponse(&p)
	}
	return result, nil
}

// GetByStudentID returns all non-archived payments of a student.
func (s *service) GetByStudentID(ctx context.Context, studentIDStr string) ([]PaymentResponse, error) {
	stID, err := primitive.ObjectIDFromHex(strings.TrimSpace(studentIDStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	payments, err := s.repo.GetByStudentID(ctx, stID)
	if err != nil {
		return nil, err
	}

	result := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		result[i] = ToPaymentResponse(&p)
	}
	return result, nil
}

// Update partially updates an existing payment.
func (s *service) Update(ctx context.Context, idStr string, req UpdatePaymentRequest) (*PaymentResponse, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil {
		if *req.Amount < 0 {
			return nil, ErrInvalidAmount
		}
		p.Amount = req.Amount
	}

	if req.DueDate != nil {
		p.DueDate = strings.TrimSpace(*req.DueDate)
	}

	if req.Status != nil {
		st := strings.TrimSpace(*req.Status)
		if !IsValidStatus(st) {
			return nil, ErrInvalidStatus
		}
		p.Status = st
	}

	if req.PaymentDate != nil {
		pDate := strings.TrimSpace(*req.PaymentDate)
		if pDate != "" && !dateRegex.MatchString(pDate) {
			return nil, ErrInvalidDateFormat
		}
		if pDate == "" {
			p.PaymentDate = nil
		} else {
			p.PaymentDate = &pDate
		}
	}

	if req.PaymentMethod != nil {
		m := strings.TrimSpace(*req.PaymentMethod)
		if !IsValidPaymentMethod(m) {
			return nil, ErrInvalidPaymentMethod
		}
		p.PaymentMethod = m
	}

	if req.TransactionReference != nil {
		p.TransactionReference = strings.TrimSpace(*req.TransactionReference)
	}

	if req.Notes != nil {
		p.Notes = strings.TrimSpace(*req.Notes)
	}

	updated, err := s.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	resp := ToPaymentResponse(updated)
	return &resp, nil
}

// Archive soft-deletes a payment record to preserve financial history.
func (s *service) Archive(ctx context.Context, idStr string) error {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.SoftDelete(ctx, id)
}

// SimulateReminder increments the reminder count and records the last reminder timestamp.
func (s *service) SimulateReminder(ctx context.Context, idStr string) (*PaymentResponse, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	nowSantiago := time.Now().In(santiagoLocation)
	timestamp := nowSantiago.Format("2006-01-02 15:04")

	p.ReminderSentCount++
	p.LastReminderAt = &timestamp

	updated, err := s.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	resp := ToPaymentResponse(updated)
	return &resp, nil
}

// SimulateReceipt marks the receipt as sent and records the emission timestamp.
func (s *service) SimulateReceipt(ctx context.Context, idStr string) (*PaymentResponse, error) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return nil, ErrInvalidID
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	nowSantiago := time.Now().In(santiagoLocation)
	timestamp := nowSantiago.Format("2006-01-02 15:04")

	p.ReceiptSent = true
	p.ReceiptSentAt = &timestamp

	updated, err := s.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	resp := ToPaymentResponse(updated)
	return &resp, nil
}

// GetFinancialSummary calculates live financial totals and transaction counts.
func (s *service) GetFinancialSummary(ctx context.Context) (*FinancialSummaryResponse, error) {
	sum, err := s.repo.GetFinancialSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &FinancialSummaryResponse{
		TotalPaidAmount:    sum.TotalPaidAmount,
		TotalPendingAmount: sum.TotalPendingAmount,
		PaidCount:          sum.PaidCount,
		PendingCount:       sum.PendingCount,
		TotalCount:         sum.TotalCount,
	}, nil
}
