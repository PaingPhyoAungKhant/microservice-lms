package usecases

import (
	"context"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type UpdateUserInput struct {
	ID       string
	Username *string
	Email    *string
	Role     *string
	Status   *string
}

type UpdateUserUseCase struct {
	userRepo  repositories.UserRepository
	publisher messaging.Publisher
	logger    *logger.Logger
}

func NewUpdateUserUseCase(
	userRepo repositories.UserRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo:  userRepo,
		publisher: publisher,
		logger:    logger,
	}
}

func (uc *UpdateUserUseCase) Execute(
	ctx context.Context,
	input UpdateUserInput,
) (*dtos.UserDTO, error) {
	user, err := uc.userRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Email != nil {
		email, err := valueobjects.NewEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	if input.Role != nil {
		role, err := valueobjects.NewRole(*input.Role)
		if err != nil {
			return nil, err
		}
		user.Role = role
	}

	if input.Status != nil {
		status, err := valueobjects.NewStatus(*input.Status)
		if err != nil {
			return nil, err
		}
		user.Status = status
	}

	user.UpdatedAt = time.Now().UTC()

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	event := events.UserUpdatedEvent{
		ID:            user.ID,
		Email:         user.Email.String(),
		Username:      user.Username,
		Role:          user.Role.String(),
		Status:        user.Status.String(),
		EmailVerified: user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(
			ctx,
			events.EventTypeUserUpdated,
			event,
		); err != nil {
			uc.logger.Error("Failed to publish user updated event", zap.Error(err))
		}
	}

	uc.logger.Info("User Updated Successfully.",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email.String()),
	)

	var dto dtos.UserDTO
	dto.FromEntity(user)
	return &dto, nil
}
