package schedule

import "context"

// StudentSummary contains the minimal student metadata required by the Schedule module.
type StudentSummary struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Archived bool   `json:"archived"`
}

// StudentReader defines the read-only interface required by Schedule to validate and display participants.
// This interface allows the Schedule module to remain completely decoupled from the Students repository,
// enabling a seamless future transition to a remote HTTP microservice without modifying domain logic.
type StudentReader interface {
	GetStudentsSummary(ctx context.Context, ids []string) ([]StudentSummary, error)
}
