package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindUserInput struct {
	SearchQuery   *string
	Role          *valueobjects.Role
	Status        *valueobjects.Status
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *repositories.SortDirection
}

type FindUserOutput struct {
	Users []dtos.UserDTO `json:"users"`
	Total int            `json:"total"`
}

type FindUserUseCase struct {
	userRepo repositories.UserRepository
	logger   *logger.Logger
}

func NewFindUserUseCase(userRepo repositories.UserRepository, logger *logger.Logger) *FindUserUseCase {
	return &FindUserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *FindUserUseCase) Execute(ctx context.Context, input FindUserInput) (*FindUserOutput, error) {
	query := repositories.UserQuery{
		SearchQuery:   input.SearchQuery,
		Role:          input.Role,
		Status:        input.Status,
		Limit:         input.Limit,
		Offset:        input.Offset,
		SortColumn:    input.SortColumn,
		SortDirection: input.SortDirection,
	}

	result, err := uc.userRepo.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dtos.UserDTO, len(result.Users))
	for i, user := range result.Users {
		userDTOs[i].FromEntity(user)
	}

	return &FindUserOutput{
		Users: userDTOs,
		Total: result.Total,
	}, nil
}