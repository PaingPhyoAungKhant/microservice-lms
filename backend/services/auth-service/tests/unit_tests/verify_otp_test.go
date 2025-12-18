package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVerifyOTP_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		OTP: "123456",
	})
	require.ErrorIs(t, err, usecases.ErrEmailRequired)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.ErrorMessage, "email is required")

	redis.AssertNotCalled(t, "GetUserFromForgotPasswordOTP", mock.Anything, mock.Anything)
}

func TestVerifyOTP_MissingOTP(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "user@example.com",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.ErrorMessage, "invalid OTP")

	redis.AssertNotCalled(t, "GetUserFromForgotPasswordOTP", mock.Anything, mock.Anything)
}

func TestVerifyOTP_InvalidOTP(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "invalid-otp").Return("", errors.New("not found")).Once()

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "user@example.com",
		OTP:   "invalid-otp",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.ErrorMessage, "invalid OTP")

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
}

func TestVerifyOTP_OTPNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "123456").Return("", nil).Once()

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "user@example.com",
		OTP:   "123456",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
}

func TestVerifyOTP_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "123456").Return("user-id", nil).Once()
	repo.On("FindByID", mock.Anything, "user-id").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "user@example.com",
		OTP:   "123456",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyOTP_EmailMismatch(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "123456").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "different@example.com",
		OTP:   "123456",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyOTP_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "123456").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	redis.On("RevokeForgotPasswordOTP", mock.Anything, "123456").Return(nil).Once()
	redis.On("StoreResetPasswordToken", mock.Anything, user.ID, mock.Anything).Return(nil).Once()

	uc := usecases.NewVerifyOTPUseCase(repo, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyOTPInput{
		Email: "user@example.com",
		OTP:   "123456",
	})
	require.NoError(t, err)
	assert.True(t, result.IsValid)
	assert.NotEmpty(t, result.PasswordResetToken)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

