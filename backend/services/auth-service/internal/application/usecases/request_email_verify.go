package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

type RequestEmailVerifyInput struct {
	Email string
}

type RequestEmailVerifyOutput struct {
	Message string
}

type RequestEmailVerifyUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	logger *logger.Logger
	redis utils.RedisInterface
	apiGatewayURL string
}

func NewRequestEmailVerifyUseCase(userRepo repositories.UserRepository, publisher messaging.Publisher, logger *logger.Logger, redis utils.RedisInterface, apiGatewayURL string) *RequestEmailVerifyUseCase {
	return &RequestEmailVerifyUseCase{
		userRepo: userRepo,
		publisher: publisher,
		logger: logger,
		redis: redis,
		apiGatewayURL: apiGatewayURL,
	}
}

func (uc *RequestEmailVerifyUseCase) Execute(ctx context.Context, input RequestEmailVerifyInput) (*RequestEmailVerifyOutput, error) {
	validationErr := uc.validateInput(input)
	if validationErr != nil {
		return nil, validationErr
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

	if user.EmailVerified {
		return &RequestEmailVerifyOutput{Message: "Email already verified"}, nil
	}

	token := uuid.New().String()
	err = uc.redis.StoreVerifyEmailToken(ctx, user.ID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to store verify email token: %w", err)
	}

	emailVerificationURL := utils.GenerateEmailVerificationURL(uc.apiGatewayURL, token)

	event := events.AuthUserRequestedEmailVerificationEvent{
		ID: user.ID,
		Email: user.Email.String(),
		EmailVerificationURL: emailVerificationURL,
	}

	if err := uc.publisher.Publish(ctx, events.EventTypeAuthUserRequestedEmailVerification, event); err != nil {
		uc.logger.Error("failed to publish email verification request event", zap.Error(err))
	}

	return &RequestEmailVerifyOutput{Message: "Verification email sent"}, nil
}

func (uc *RequestEmailVerifyUseCase) validateInput(input RequestEmailVerifyInput) error {
	if input.Email == "" {
		return errors.New("email is required")
	}
	return nil
}