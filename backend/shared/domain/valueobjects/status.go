package valueobjects

import (
	"errors"
)

var ErrInvalidStatus = errors.New("invalid status")

type Status string

const (
	StatusActive   Status = "active"
	StatusInActive Status = "inactive"
	StatusPending  Status = "pending"
	StatusBanned   Status = "banned"
)

func NewStatus(status string) (Status, error) {
	s := Status(status)
	switch s {
	case StatusActive, StatusInActive, StatusPending, StatusBanned:
		return s, nil
	default:
		return "", ErrInvalidStatus
	}
}

func (s Status) String() string {
	return string(s)
}

func (s Status) IsActive() bool {
	return s == StatusActive
}

func (s Status) CanLogin() bool {
	return s == StatusActive
}

func (s Status) IsBanned() bool {
	return s == StatusBanned
}

