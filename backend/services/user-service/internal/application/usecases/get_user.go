package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID")
)

type GetUserInput struct {
	UserID string
}

type GetUserUseCase struct {
	userRepo repositories.UserRepository
	logger   *logger.Logger
}

func NewGetUserUseCase(userRepo repositories.UserRepository, logger *logger.Logger) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*dtos.UserDTO, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	var dto dtos.UserDTO
	dto.FromEntity(user)
	return &dto, nil
}
