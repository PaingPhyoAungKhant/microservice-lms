package usecases

import (
	"context"
	"errors"

	authUtils "github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/utils"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP")
)

type VerifyOTPInput struct {
	Email string
	OTP string
}

type VerifyOTPOutput struct {
	IsValid            bool   `json:"is_valid"`
	PasswordResetToken string `json:"password_reset_token"`
	ErrorMessage       string `json:"error_message"`
}

type VerifyOTPUseCase struct {
	userRepo repositories.UserRepository
	logger *logger.Logger
	redis utils.RedisInterface

}

func NewVerifyOTPUseCase(userRepo repositories.UserRepository, logger *logger.Logger, redis utils.RedisInterface) *VerifyOTPUseCase {
	return &VerifyOTPUseCase{userRepo: userRepo, logger: logger, redis: redis}
}

func (uc *VerifyOTPUseCase) Execute(ctx context.Context, input VerifyOTPInput) (*VerifyOTPOutput, error) {
	output := &VerifyOTPOutput{IsValid: false}

	if err := uc.validateInput(input); err != nil {
		output.ErrorMessage = err.Error()
		return output, err
	}

	userID, err := uc.redis.GetUserFromForgotPasswordOTP(ctx, input.OTP)
	if err != nil {
		uc.logger.Error("failed to get user from OTP", zap.Error(err))
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	if userID == "" {
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user", zap.Error(err))
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	if user == nil {
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	if user.Email.String() != input.Email {
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	if err := uc.redis.RevokeForgotPasswordOTP(ctx, input.OTP); err != nil {
		uc.logger.Error("failed to revoke OTP", zap.Error(err))
		output.ErrorMessage = ErrInvalidOTP.Error()
		return output, ErrInvalidOTP
	}

	passwordResetToken, err := authUtils.GeneratePasswordResetToken()
	if err != nil {
		uc.logger.Error("failed to generate password reset token", zap.Error(err))
		output.ErrorMessage = ErrInternalServerError.Error()
		return output, ErrInternalServerError
	}

	if err := uc.redis.StoreResetPasswordToken(ctx, userID, passwordResetToken); err != nil {
		uc.logger.Error("failed to store reset password token", zap.Error(err))
		output.ErrorMessage = ErrInternalServerError.Error()
		return output, ErrInternalServerError
	}

	output.IsValid = true
	output.PasswordResetToken = passwordResetToken
	return output, nil
}

func (uc *VerifyOTPUseCase) validateInput(input VerifyOTPInput) error {
	if input.Email == "" {
		return ErrEmailRequired
	}
	if input.OTP == "" {
		return ErrInvalidOTP
	}
	return nil
}