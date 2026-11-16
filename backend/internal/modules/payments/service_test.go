package payments

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockRepo is an in-memory implementation of payments.Repository for isolated unit testing.
type mockRepo struct {
	payments map[primitive.ObjectID]*Payment
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		payments: make(map[primitive.ObjectID]*Payment),
	}
}

func (m *mockRepo) EnsureIndexes(ctx context.Context) error {
	return nil
}

func (m *mockRepo) Create(ctx context.Context, p *Payment) (*Payment, error) {
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	m.payments[p.ID] = p
	return p, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*Payment, error) {
	p, ok := m.payments[id]
	if !ok || p.ArchivedAt != nil {
		return nil, ErrPaymentNotFound
	}
	return p, nil
}

func (m *mockRepo) List(ctx context.Context, filter PaymentFilter) ([]Payment, error) {
	var list []Payment
	for _, p := range m.payments {
		if p.ArchivedAt != nil {
			continue
		}
		if filter.StudentID != "" && p.StudentID.Hex() != filter.StudentID {
			continue
		}
		if filter.Status != "" && p.Status != filter.Status {
			continue
		}
		list = append(list, *p)
	}
	return list, nil
}

func (m *mockRepo) GetByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]Payment, error) {
	var list []Payment
	for _, p := range m.payments {
		if p.ArchivedAt == nil && p.StudentID == studentID {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *mockRepo) GetPendingByStudentID(ctx context.Context, studentID primitive.ObjectID) (*Payment, error) {
	for _, p := range m.payments {
		if p.ArchivedAt == nil && p.StudentID == studentID && p.Status == StatusPending {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) Update(ctx context.Context, p *Payment) (*Payment, error) {
	existing, ok := m.payments[p.ID]
	if !ok || existing.ArchivedAt != nil {
		return nil, ErrPaymentNotFound
	}
	p.UpdatedAt = time.Now().UTC()
	m.payments[p.ID] = p
	return p, nil
}

func (m *mockRepo) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	p, ok := m.payments[id]
	if !ok || p.ArchivedAt != nil {
		return ErrPaymentNotFound
	}
	now := time.Now().UTC()
	p.ArchivedAt = &now
	p.UpdatedAt = now
	return nil
}

func (m *mockRepo) GetFinancialSummary(ctx context.Context) (*FinancialSummary, error) {
	summary := &FinancialSummary{}
	for _, p := range m.payments {
		if p.ArchivedAt != nil {
			continue
		}
		summary.TotalCount++
		if p.Status == StatusPaid {
			summary.PaidCount++
			if p.Amount != nil {
				summary.TotalPaidAmount += *p.Amount
			}
		} else if p.Status == StatusPending {
			summary.PendingCount++
			if p.Amount != nil {
				summary.TotalPendingAmount += *p.Amount
			}
		}
	}
	return summary, nil
}

// mockStudentReader simulates StudentReader.
type mockStudentReader struct {
	students        map[string]*StudentInfo
	activatedCalled bool
}

func newMockStudentReader() *mockStudentReader {
	return &mockStudentReader{
		students: make(map[string]*StudentInfo),
	}
}

func (m *mockStudentReader) GetStudent(ctx context.Context, id string) (*StudentInfo, error) {
	s, ok := m.students[id]
	if !ok || s.Archived {
		return nil, ErrStudentNotFound
	}
	return s, nil
}

func (m *mockStudentReader) ActivateMembershipIfReviewing(ctx context.Context, studentID string) error {
	if s, ok := m.students[studentID]; ok && s.MembershipStatus == "Por revisar" {
		s.MembershipStatus = "Activo"
		m.activatedCalled = true
	}
	return nil
}

func TestRegisterNewPayment(t *testing.T) {
	repo := newMockRepo()
	reader := newMockStudentReader()
	studentID := primitive.NewObjectID().Hex()
	reader.students[studentID] = &StudentInfo{
		ID:               studentID,
		FullName:         "Juan Pérez",
		PlanName:         "Plan 3 Veces",
		Frequency:        "3 veces/sem.",
		MembershipStatus: "Por revisar",
	}

	svc := NewService(repo, reader)
	amt := int64(45000)

	resp, err := svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID:     studentID,
		Amount:        &amt,
		PaymentMethod: MethodTransfer,
		PaymentDate:   "2026-11-06",
		Notes:         "Comprobante recibido",
	})

	if err != nil {
		t.Fatalf("unexpected error registering payment: %v", err)
	}
	if resp.Amount == nil || *resp.Amount != 45000 {
		t.Errorf("expected amount 45000, got %v", resp.Amount)
	}
	if resp.Status != StatusPaid {
		t.Errorf("expected status 'Pagado', got %s", resp.Status)
	}
	if !reader.activatedCalled {
		t.Errorf("expected student membership to be activated from 'Por revisar'")
	}
}

func TestRegisterPaymentUpdatesPending(t *testing.T) {
	repo := newMockRepo()
	reader := newMockStudentReader()
	stOID := primitive.NewObjectID()
	studentID := stOID.Hex()
	reader.students[studentID] = &StudentInfo{
		ID:               studentID,
		FullName:         "María López",
		PlanName:         "Plan 2 Veces",
		Frequency:        "2 veces/sem.",
		MembershipStatus: "Activo",
	}

	// Seed existing pending payment
	initialPending := &Payment{
		StudentID:   stOID,
		StudentName: "María López",
		PlanName:    "Plan 2 Veces",
		Frequency:   "2 veces/sem.",
		Status:      StatusPending,
		DueDate:     "2026-11-15",
	}
	_, _ = repo.Create(context.Background(), initialPending)

	svc := NewService(repo, reader)
	amt := int64(38000)

	resp, err := svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID:     studentID,
		Amount:        &amt,
		PaymentMethod: MethodDebit,
		PaymentDate:   "2026-11-06",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != initialPending.ID.Hex() {
		t.Errorf("expected to update pending payment id %s, got %s", initialPending.ID.Hex(), resp.ID)
	}
	if resp.Status != StatusPaid {
		t.Errorf("expected status to be updated to 'Pagado', got %s", resp.Status)
	}
}

func TestValidationErrors(t *testing.T) {
	repo := newMockRepo()
	reader := newMockStudentReader()
	studentID := primitive.NewObjectID().Hex()
	reader.students[studentID] = &StudentInfo{ID: studentID, FullName: "Test"}

	svc := NewService(repo, reader)

	// Negative amount
	neg := int64(-500)
	_, err := svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID: studentID,
		Amount:    &neg,
	})
	if err != ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

	// Invalid date
	amt := int64(10000)
	_, err = svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID:   studentID,
		Amount:      &amt,
		PaymentDate: "06-11-2026", // Wrong format
	})
	if err != ErrInvalidDateFormat {
		t.Errorf("expected ErrInvalidDateFormat, got %v", err)
	}

	// Invalid method
	_, err = svc.RegisterPayment(context.Background(), RegisterPaymentRequest{
		StudentID:     studentID,
		Amount:        &amt,
		PaymentMethod: "Bitcoins", // Unsupported
	})
	if err != ErrInvalidPaymentMethod {
		t.Errorf("expected ErrInvalidPaymentMethod, got %v", err)
	}
}

func TestFinancialSummary(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, nil)

	p1Amt := int64(45000)
	p2Amt := int64(30000)
	_, _ = repo.Create(context.Background(), &Payment{
		Status: StatusPaid,
		Amount: &p1Amt,
	})
	_, _ = repo.Create(context.Background(), &Payment{
		Status: StatusPending,
		Amount: &p2Amt,
	})

	summary, err := svc.GetFinancialSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.PaidCount != 1 || summary.PendingCount != 1 || summary.TotalCount != 2 {
		t.Errorf("unexpected counts in summary: %+v", summary)
	}
	if summary.TotalPaidAmount != 45000 || summary.TotalPendingAmount != 30000 {
		t.Errorf("unexpected amounts in summary: %+v", summary)
	}
}

func TestSanitizeCSVCell(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"=cmd|' /C calc'!A0", "'=cmd|' /C calc'!A0"},
		{"+12345", "'+12345"},
		{"-formula", "'-formula"},
		{"@SUM(A1:A10)", "'@SUM(A1:A10)"},
		{"Normal text", "Normal text"},
		{"Carlos Mendoza", "Carlos Mendoza"},
	}

	for _, tc := range testCases {
		res := SanitizeCSVCell(tc.input)
		if res != tc.expected {
			t.Errorf("SanitizeCSVCell(%q) = %q; want %q", tc.input, res, tc.expected)
		}
	}
}

func TestSimulateReminderAndReceipt(t *testing.T) {
	repo := newMockRepo()
	p, _ := repo.Create(context.Background(), &Payment{
		Status: StatusPending,
	})
	svc := NewService(repo, nil)

	remResp, err := svc.SimulateReminder(context.Background(), p.ID.Hex())
	if err != nil {
		t.Fatalf("unexpected reminder error: %v", err)
	}
	if remResp.ReminderSentCount != 1 || remResp.LastReminderAt == nil {
		t.Errorf("expected reminder count 1 and timestamp, got %+v", remResp)
	}

	recResp, err := svc.SimulateReceipt(context.Background(), p.ID.Hex())
	if err != nil {
		t.Fatalf("unexpected receipt error: %v", err)
	}
	if !recResp.ReceiptSent || recResp.ReceiptSentAt == nil {
		t.Errorf("expected receiptSent true and timestamp, got %+v", recResp)
	}
}
