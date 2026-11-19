package payments

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Allowed payment statuses
const (
	StatusPaid      = "Pagado"
	StatusPending   = "Pendiente"
	StatusUndefined = "Por definir"
)

// Allowed payment methods
const (
	MethodSumUpLink    = "SumUp (Link de pago)"
	MethodSumUpPOS     = "SumUp (POS Tarjeta)"
	MethodTransferBank = "Transferencia Bancaria"
	MethodCash         = "Efectivo"
	MethodTransfer     = "Transferencia"
	MethodWebpay       = "Webpay"
	MethodServipag     = "Servipag"
	MethodCredit       = "Crédito"
	MethodDebit        = "Débito"
	MethodUndefined    = "Por definir"
)

// IsValidStatus checks whether a status string is valid.
func IsValidStatus(status string) bool {
	switch status {
	case StatusPaid, StatusPending, StatusUndefined:
		return true
	default:
		return false
	}
}

// IsValidPaymentMethod checks whether a payment method string is supported.
func IsValidPaymentMethod(method string) bool {
	switch method {
	case MethodSumUpLink, MethodSumUpPOS, MethodTransferBank, MethodCash, MethodTransfer, MethodWebpay, MethodServipag, MethodCredit, MethodDebit, MethodUndefined:
		return true
	default:
		return false
	}
}

// SanitizeCSVCell prevents CSV Formula Injection (DDE injection).
// If a string starts with '=', '+', '-', '@', '\t' or '\r', it prepends a single quote.
func SanitizeCSVCell(val string) string {
	trimmed := strings.TrimSpace(val)
	if len(trimmed) > 0 {
		first := trimmed[0]
		if first == '=' || first == '+' || first == '-' || first == '@' || first == '\t' || first == '\r' {
			return "'" + val
		}
	}
	return val
}

// Payment represents a financial transaction or pending billing record in INVERFIT.
type Payment struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	StudentID            primitive.ObjectID `bson:"student_id"`
	StudentName          string             `bson:"student_name"`
	PlanName             string             `bson:"plan_name"`
	Frequency            string             `bson:"frequency"`
	DueDate              string             `bson:"due_date"`             // "YYYY-MM-DD" or "Por definir"
	Amount               *int64             `bson:"amount,omitempty"`    // Integer CLP (null if undefined)
	Status               string             `bson:"status"`               // "Pagado" | "Pendiente" | "Por definir"
	PaymentDate          *string            `bson:"payment_date,omitempty"` // "YYYY-MM-DD" or nil
	PaymentMethod        string             `bson:"payment_method"`        // "Transferencia", "Webpay", etc.
	TransactionReference string             `bson:"transaction_reference,omitempty"`
	ReminderSentCount    int                `bson:"reminder_sent_count"`
	LastReminderAt       *string            `bson:"last_reminder_at,omitempty"`
	ReceiptSent          bool               `bson:"receipt_sent"`
	ReceiptSentAt        *string            `bson:"receipt_sent_at,omitempty"`
	Notes                string             `bson:"notes,omitempty"`
	ArchivedAt           *time.Time         `bson:"archived_at,omitempty"` // Soft delete
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
}

// FinancialSummary holds aggregated financial KPIs calculated from persisted payment records.
type FinancialSummary struct {
	TotalPaidAmount    int64 `json:"totalPaidAmount"`
	TotalPendingAmount int64 `json:"totalPendingAmount"`
	PaidCount          int   `json:"paidCount"`
	PendingCount       int   `json:"pendingCount"`
	TotalCount         int   `json:"totalCount"`
}
