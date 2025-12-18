package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
)

type VerifyEmailInput struct {
	Email string 
	Token string
}

type VerifyEmailOutput struct {
	Message string
}

type VerifyEmailUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	logger *logger.Logger
	redis utils.RedisInterface
}

func NewVerifyEmailUseCase(userRepo repositories.UserRepository, publisher messaging.Publisher, logger *logger.Logger, redis utils.RedisInterface) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepo: userRepo,
		publisher: publisher,
		logger: logger,
		redis: redis,
	}
}

func (uc *VerifyEmailUseCase) Execute(ctx context.Context, input VerifyEmailInput) (*VerifyEmailOutput, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	userID, err := uc.redis.GetUserFromVerifyEmailToken(ctx, input.Token)
	if err != nil {
		return nil, errors.New("invalid or expired verification token")
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if input.Email != "" {
		email, err := valueobjects.NewEmail(input.Email)
		if err != nil {
			return nil, err
		}
		if user.Email.String() != email.String() {
			return nil, errors.New("email does not match token")
		}
	}

	if user.EmailVerified {
		return &VerifyEmailOutput{Message: "Email already verified"}, nil
	}

	if err := uc.userRepo.UpdateEmailVerified(ctx, userID, true); err != nil {
		return nil, err
	}

	uc.redis.RevokeVerifyEmailToken(ctx, input.Token)

	return &VerifyEmailOutput{Message: "Email verified successfully"}, nil
}

func (uc *VerifyEmailUseCase) validateInput(input VerifyEmailInput) error {
	if input.Token == "" {
		return errors.New("token is required")
	}
	return nil
}