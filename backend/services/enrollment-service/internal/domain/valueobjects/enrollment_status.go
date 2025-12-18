package valueobjects

import "errors"

type EnrollmentStatus string

const (
	EnrollmentStatusPending   EnrollmentStatus = "pending"
	EnrollmentStatusApproved  EnrollmentStatus = "approved"
	EnrollmentStatusRejected  EnrollmentStatus = "rejected"
	EnrollmentStatusCompleted EnrollmentStatus = "completed"
)

var (
	ErrInvalidEnrollmentStatus = errors.New("invalid enrollment status")
)

func NewEnrollmentStatus(status string) (EnrollmentStatus, error) {
	switch status {
	case "pending", "approved", "rejected", "completed":
		return EnrollmentStatus(status), nil
	default:
		return "", ErrInvalidEnrollmentStatus
	}
}

func (s EnrollmentStatus) String() string {
	return string(s)
}

func (s EnrollmentStatus) IsValid() bool {
	return s == EnrollmentStatusPending ||
		s == EnrollmentStatusApproved ||
		s == EnrollmentStatusRejected ||
		s == EnrollmentStatusCompleted
}

