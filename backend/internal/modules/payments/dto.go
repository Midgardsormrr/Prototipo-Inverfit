package payments

import "time"

// RegisterPaymentRequest represents the payload received when creating or confirming a payment.
type RegisterPaymentRequest struct {
	StudentID            string `json:"studentId" binding:"required"`
	Amount               *int64 `json:"amount"`
	DueDate              string `json:"dueDate"`
	PaymentDate          string `json:"paymentDate"`
	PaymentMethod        string `json:"paymentMethod"`
	TransactionReference string `json:"transactionReference"`
	Notes                string `json:"notes"`
}

// UpdatePaymentRequest represents the payload for updating an existing payment.
type UpdatePaymentRequest struct {
	Amount               *int64  `json:"amount"`
	DueDate              *string `json:"dueDate"`
	Status               *string `json:"status"`
	PaymentDate          *string `json:"paymentDate"`
	PaymentMethod        *string `json:"paymentMethod"`
	TransactionReference *string `json:"transactionReference"`
	Notes                *string `json:"notes"`
}

// PaymentResponse represents the payment object exposed to the Vue 3 frontend.
// It is 100% compatible with the TypeScript 'Payment' interface.
type PaymentResponse struct {
	ID                   string  `json:"id"`
	StudentID            string  `json:"studentId"`
	StudentName          string  `json:"studentName"`
	PlanName             string  `json:"planName"`
	Frequency            string  `json:"frequency"`
	DueDate              string  `json:"dueDate"`
	Amount               *int64  `json:"amount"`
	Status               string  `json:"status"`
	PaymentDate          *string `json:"paymentDate"`
	PaymentMethod        string  `json:"paymentMethod"`
	TransactionReference string  `json:"transactionReference,omitempty"`
	ReminderSentCount    int     `json:"reminderSentCount"`
	LastReminderAt       *string `json:"lastReminderAt"`
	ReceiptSent          bool    `json:"receiptSent"`
	ReceiptSentAt        *string `json:"receiptSentAt"`
	Notes                string  `json:"notes,omitempty"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

// FinancialSummaryResponse represents aggregated financial metrics.
type FinancialSummaryResponse struct {
	TotalPaidAmount    int64 `json:"totalPaidAmount"`
	TotalPendingAmount int64 `json:"totalPendingAmount"`
	PaidCount          int   `json:"paidCount"`
	PendingCount       int   `json:"pendingCount"`
	TotalCount         int   `json:"totalCount"`
}

// PaymentFilter defines query filter parameters for listing payments.
type PaymentFilter struct {
	StudentID string
	Status    string
	Search    string
	From      string // YYYY-MM-DD
	To        string // YYYY-MM-DD
}

// ToPaymentResponse transforms a domain Payment model into a JSON-friendly DTO.
func ToPaymentResponse(p *Payment) PaymentResponse {
	return PaymentResponse{
		ID:                   p.ID.Hex(),
		StudentID:            p.StudentID.Hex(),
		StudentName:          p.StudentName,
		PlanName:             p.PlanName,
		Frequency:            p.Frequency,
		DueDate:              p.DueDate,
		Amount:               p.Amount,
		Status:               p.Status,
		PaymentDate:          p.PaymentDate,
		PaymentMethod:        p.PaymentMethod,
		TransactionReference: p.TransactionReference,
		ReminderSentCount:    p.ReminderSentCount,
		LastReminderAt:       p.LastReminderAt,
		ReceiptSent:          p.ReceiptSent,
		ReceiptSentAt:        p.ReceiptSentAt,
		Notes:                p.Notes,
		CreatedAt:            p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            p.UpdatedAt.Format(time.RFC3339),
	}
}
