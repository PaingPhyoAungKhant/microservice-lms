package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
)

var (
	ErrInvalidPasswordResetToken = errors.New("invalid password reset token")
	ErrInvalidNewPassword = errors.New("invalid new password")
	ErrInvalidIPAddress = errors.New("invalid IP address")
	ErrInvalidUserAgent = errors.New("invalid user agent")
)

type ResetPasswordInput struct {
	Token string
	NewPassword string
	IPAddress string
	UserAgent string
}

type ResetPasswordOutput struct {
	Message string
}

type ResetPasswordUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	logger *logger.Logger
	redis utils.RedisInterface
}

func NewResetPasswordUseCase(userRepo repositories.UserRepository, publisher messaging.Publisher, logger *logger.Logger, redis utils.RedisInterface) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{userRepo: userRepo, publisher: publisher, logger: logger, redis: redis}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, error) {
	output := &ResetPasswordOutput{Message: "Password reset successful"}

	if err := uc.validateInput(input); err != nil {
		output.Message = err.Error()
		return output, err
	}

	if err := utils.ValidatePassword(input.NewPassword); err != nil {
		output.Message = err.Error()
		return output, err
	}

	passwordHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		output.Message = err.Error()
		return output, err
	}

	userID, err := uc.redis.GetUserFromResetPasswordToken(ctx, input.Token)
	if err != nil {
		output.Message = err.Error()
		return output, err
	}

	if userID == "" {
		output.Message = ErrInvalidPasswordResetToken.Error()
		return output, ErrInvalidPasswordResetToken
	}
	
	if err := uc.redis.RevokeResetPasswordToken(ctx, input.Token); err != nil {
		output.Message = err.Error()
		return output, err
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		output.Message = err.Error()
		return output, err
	}

	if user == nil {
		output.Message = ErrUserNotFound.Error()
		return output, ErrUserNotFound
	}
	
	if err := uc.userRepo.UpdatePassword(ctx, userID, passwordHash); err != nil {
		output.Message = err.Error()
		return output, err
	}
	event := events.AuthUserResetPasswordEvent{
		ID: userID,
		Username: user.Username,
		Email: user.Email.String(),
		Role: user.Role.String(),
		Status: user.Status.String(),
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		PublishedAt: time.Now(),
	}
	if err := uc.publisher.Publish(ctx, events.EventTypeAuthUserResetPassword, event); err != nil {
		output.Message = err.Error()
		return output, err
	}
	return output, nil
}

func (uc *ResetPasswordUseCase) validateInput(input ResetPasswordInput) error {
	if input.Token == "" {
		return ErrInvalidPasswordResetToken
	}
	if input.NewPassword == "" {
		return ErrInvalidPassword
	}
	if input.IPAddress == "" {
		return ErrInvalidIPAddress
	}
	if input.UserAgent == "" {
		return ErrInvalidUserAgent
	}
	return nil
}