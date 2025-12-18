package events

import "time"

type AuthStudentRegisteredEvent struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email"`
	Username            string     `json:"username"`
	Role                string     `json:"role"`
	Status              string     `json:"status"`
	EmailVerified       bool       `json:"email_verified"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	EmailVerificationURL string    `json:"email_verification_url"`
}

type AuthUserLoggedInEvent struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AuthUserForgotPasswordEvent struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	PublishedAt time.Time `json:"publishedAt"`
}

type AuthUserResetPasswordEvent struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	IPAddress   string    `json:"ipAddress"`
	UserAgent   string    `json:"userAgent"`
	PublishedAt time.Time `json:"publishedAt"`
}

type AuthUserLoggedOutEvent struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AuthUserRequestedEmailVerificationEvent struct {
	ID                  string `json:"id"`
	Email               string `json:"email"`
	EmailVerificationURL string `json:"email_verification_url"`
}

