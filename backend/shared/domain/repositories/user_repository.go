// Package repositories
package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
)

type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

type UserQuery struct {
	SearchQuery   *string
	Role          *valueobjects.Role
	Status        *valueobjects.Status
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *SortDirection
}

type UserQueryResult struct {
	Users []*entities.User
	Total int
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateEmailVerified(ctx context.Context, userID string, verified bool) error
	Find(ctx context.Context, query UserQuery) (*UserQueryResult, error)
}

