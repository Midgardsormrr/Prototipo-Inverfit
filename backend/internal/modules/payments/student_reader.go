package payments

import "context"

// StudentInfo contains student metadata needed by the Payments module.
type StudentInfo struct {
	ID               string `json:"id"`
	FullName         string `json:"fullName"`
	PlanName         string `json:"planName"`
	Frequency        string `json:"frequency"`
	MembershipStatus string `json:"membershipStatus"`
	Archived         bool   `json:"archived"`
}

// StudentReader defines the decoupled contract for querying and updating student status.
type StudentReader interface {
	GetStudent(ctx context.Context, id string) (*StudentInfo, error)
	ActivateMembershipIfReviewing(ctx context.Context, studentID string) error
}
