// Package dtos
package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
)

type UserDTO struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Username       string     `json:"username"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	EmailVerified  bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (d *UserDTO) FromEntity(user *entities.User) {
	d.ID = user.ID
	d.Email = user.Email.String()
	d.Username = user.Username
	d.Role = user.Role.String()
	d.Status = user.Status.String()
	d.EmailVerified = user.EmailVerified
	d.EmailVerifiedAt = user.EmailVerifiedAt
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
}

