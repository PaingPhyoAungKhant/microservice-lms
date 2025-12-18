package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
)

var (
	ErrCannotBanAdmin       = errors.New("cannot ban admin user")
	ErrCannotActivateBanned = errors.New("cannot activate banned user")
)

type User struct {
	ID            string
	Email         valueobjects.Email
	Username      string
	PasswordHash  string
	Role          valueobjects.Role
	Status        valueobjects.Status
	EmailVerified bool
	EmailVerifiedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUser(email valueobjects.Email, username string, role valueobjects.Role, password_hash string) *User {
	now := time.Now().UTC()
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		Username:     username,
		PasswordHash: password_hash,
		Role:         role,
		Status:       valueobjects.StatusActive,
		EmailVerified: false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) ChangePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) ChangeRole(newRole valueobjects.Role) {
	u.Role = newRole
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Activate() error {
	if u.Status.IsBanned() {
		return ErrCannotActivateBanned
	}
	u.Status = valueobjects.StatusActive
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) Deactivate() {
	u.Status = valueobjects.StatusInActive
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Ban() error {
	if u.Role.IsAdmin() {
		return ErrCannotBanAdmin
	}
	u.Status = valueobjects.StatusBanned
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) CanEnroll() bool {
	return u.Status.IsActive() && u.Role.IsStudent()
}

func (u *User) CanTeach() bool {
	return u.Status.IsActive() && u.Role.IsInstructor()
}

func (u *User) IsActive() bool {
	return u.Status.IsActive()
}

func (u *User) VerifyEmail() {
	now := time.Now().UTC()
	u.EmailVerified = true
	u.EmailVerifiedAt = &now
	u.UpdatedAt = now
}

