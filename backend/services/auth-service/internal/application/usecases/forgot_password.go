package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var (
	ErrInternalServerError = errors.New("internal server error")
)

type ForgotPasswordInput struct {
	Email string
}

type ForgotPasswordOutput struct {
	Message string
}

type ForgotPasswordUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	logger *logger.Logger

}

func NewForgotPasswordUseCase(
	userRepo repositories.UserRepository,
	 publisher messaging.Publisher,
	 logger *logger.Logger,

) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		userRepo: userRepo,
		publisher: publisher,
		logger: logger,
	}
}

func (uc *ForgotPasswordUseCase) Execute(
	ctx context.Context, 
	input ForgotPasswordInput,
) (*ForgotPasswordOutput, error) {
	
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}


	event := events.AuthUserForgotPasswordEvent{
	  ID: user.ID,
	  Email: user.Email.String(),
		Role: user.Role.String(),
		Status: user.Status.String(),
		PublishedAt: time.Now(),
	}

	if err := uc.publisher.Publish(ctx, events.EventTypeAuthUserForgotPassword, event); err != nil {
		uc.logger.Error("failed to publish forgot password event", zap.Error(err))
		return nil, ErrInternalServerError
	}

	return &ForgotPasswordOutput{
		Message: "Password reset email sent",
	}, nil
}

func (uc *ForgotPasswordUseCase) validateInput(input ForgotPasswordInput) error {
	if input.Email == "" {
		return ErrEmailRequired
	}
	return nil
}