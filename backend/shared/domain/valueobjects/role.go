// Package valueobjects
package valueobjects

import (
	"errors"
)

var ErrInvalidRole = errors.New("invalid role")

type Role string

const (
	RoleStudent    Role = "student"
	RoleInstructor Role = "instructor"
	RoleAdmin      Role = "admin"
)

func NewRole(role string) (Role, error) {
	r := Role(role)

	switch r {
	case RoleStudent, RoleInstructor, RoleAdmin:
		return r, nil
	default:
		return "", ErrInvalidRole
	}
}

func (r Role) String() string {
	return string(r)
}

func (r Role) CanCreateCourse() bool {
	return r == RoleInstructor || r == RoleAdmin
}

func (r Role) IsStudent() bool {
	return r == RoleStudent
}

func (r Role) IsInstructor() bool {
	return r == RoleInstructor
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

