// Package usecases
package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var ErrCannotDeleteAdminUser = errors.New("cannot delete admin user")

type DeleteUserInput struct {
	UserID string
}

type DeleteUserOutput struct {
	Message string `json:"message"`
}

type DeleteUserUsecase struct {
	userRepo  repositories.UserRepository
	publisher messaging.Publisher
	logger    *logger.Logger
}

func NewDeleteUserUseCase(
	userRepo repositories.UserRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteUserUsecase {
	return &DeleteUserUsecase{
		userRepo:  userRepo,
		publisher: publisher,
		logger:    logger,
	}
}

func (uc *DeleteUserUsecase) Execute(
	ctx context.Context,
	input DeleteUserInput,
) (*DeleteUserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if user.Role.IsAdmin() {
		return nil, ErrCannotDeleteAdminUser
	}

	if err := uc.userRepo.Delete(ctx, input.UserID); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	event := events.UserDeletedEvent{
		ID:        user.ID,
		DeletedAt: time.Now(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(
			ctx,
			events.EventTypeUserDeleted,
			event,
		); err != nil {
			uc.logger.Error("failed to publish user deleted event.", zap.Error(err))
		}
	}

	uc.logger.Info("User deleted Successfully",
		zap.String("user_id", user.ID),
	)

	return &DeleteUserOutput{
		Message: "User Deleted Successfully.",
	}, nil
}
