package payments

import (
	"context"
	"os"
	"strings"
	"testing"

	"inverfit-backend/internal/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepositoryNilDatabase(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	p := &Payment{}
	if _, err := repo.Create(ctx, p); err == nil {
		t.Errorf("expected error on Create with nil db")
	}
	if _, err := repo.GetByID(ctx, primitive.NewObjectID()); err == nil {
		t.Errorf("expected error on GetByID with nil db")
	}
	if _, err := repo.List(ctx, PaymentFilter{}); err == nil {
		t.Errorf("expected error on List with nil db")
	}
	if _, err := repo.Update(ctx, p); err == nil {
		t.Errorf("expected error on Update with nil db")
	}
	if err := repo.SoftDelete(ctx, primitive.NewObjectID()); err == nil {
		t.Errorf("expected error on SoftDelete with nil db")
	}
	if _, err := repo.GetFinancialSummary(ctx); err == nil {
		t.Errorf("expected error on GetFinancialSummary with nil db")
	}
}

func getPaymentsTestMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	candidates := []string{"../../.env", "../../../.env", ".env", "backend/.env"}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "MONGODB_URI=") {
				return strings.TrimSpace(strings.TrimPrefix(line, "MONGODB_URI="))
			}
		}
	}
	return ""
}

func TestRepositoryIntegrationWithMongoDB(t *testing.T) {
	uri := getPaymentsTestMongoURI()
	if uri == "" {
		t.Skip("MONGODB_URI not available; skipping live MongoDB repository test")
	}

	// Strictly connect to inverfit_test to never touch development data
	mongoDB, err := database.Connect(context.Background(), uri, "inverfit_test")
	if err != nil {
		t.Skipf("cannot connect to MongoDB: %v; skipping live test", err)
	}
	defer mongoDB.Disconnect(context.Background())

	db := mongoDB.Database
	repo := NewRepository(db)
	ctx := context.Background()

	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("failed EnsureIndexes: %v", err)
	}

	// 1. Create a payment
	amt := int64(50000)
	stID := primitive.NewObjectID()
	payDate := "2026-11-06"
	newPay := &Payment{
		StudentID:     stID,
		StudentName:   "Test Alumno Repo",
		PlanName:      "Plan Repo",
		Frequency:     "3 veces/sem.",
		DueDate:       "2026-11-30",
		Amount:        &amt,
		Status:        StatusPaid,
		PaymentDate:   &payDate,
		PaymentMethod: MethodTransfer,
		Notes:         "Test repo",
	}

	created, err := repo.Create(ctx, newPay)
	if err != nil {
		t.Fatalf("failed Create: %v", err)
	}
	if created.ID.IsZero() {
		t.Errorf("expected non-zero ObjectID")
	}

	// 2. GetByID
	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed GetByID: %v", err)
	}
	if fetched.StudentName != "Test Alumno Repo" || *fetched.Amount != 50000 {
		t.Errorf("unexpected fetched data: %+v", fetched)
	}

	// 3. List
	list, err := repo.List(ctx, PaymentFilter{StudentID: stID.Hex()})
	if err != nil {
		t.Fatalf("failed List: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("expected at least 1 payment in list")
	}

	// 4. Update
	newAmt := int64(55000)
	fetched.Amount = &newAmt
	updated, err := repo.Update(ctx, fetched)
	if err != nil {
		t.Fatalf("failed Update: %v", err)
	}
	if *updated.Amount != 55000 {
		t.Errorf("expected updated amount 55000, got %d", *updated.Amount)
	}

	// 5. Financial Summary
	summary, err := repo.GetFinancialSummary(ctx)
	if err != nil {
		t.Fatalf("failed GetFinancialSummary: %v", err)
	}
	if summary.TotalCount == 0 {
		t.Errorf("expected total count > 0")
	}

	// 6. Soft Delete
	if err := repo.SoftDelete(ctx, created.ID); err != nil {
		t.Fatalf("failed SoftDelete: %v", err)
	}

	// 7. Verify GetByID returns ErrPaymentNotFound after soft delete
	_, err = repo.GetByID(ctx, created.ID)
	if err != ErrPaymentNotFound {
		t.Errorf("expected ErrPaymentNotFound after soft delete, got %v", err)
	}
}
