package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVerifyOTP_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("otp@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "otpuser", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	otp := "123456"
	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, otp).Return(user.ID, nil).Once()
	redis.On("RevokeForgotPasswordOTP", mock.Anything, otp).Return(nil).Once()
	redis.On("StoreResetPasswordToken", mock.Anything, user.ID, mock.Anything).Return(nil).Once()

	verifyOTPUC := usecases.NewVerifyOTPUseCase(userRepo, logger, redis)

	input := usecases.VerifyOTPInput{
		Email: "otp@example.com",
		OTP:   otp,
	}

	result, err := verifyOTPUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.True(t, result.IsValid)
	assert.NotEmpty(t, result.PasswordResetToken)

	redis.AssertExpectations(t)
}

func TestVerifyOTP_Integration_InvalidOTP(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, "invalid-otp").Return("", nil).Once()

	verifyOTPUC := usecases.NewVerifyOTPUseCase(userRepo, logger, redis)

	input := usecases.VerifyOTPInput{
		Email: "otp@example.com",
		OTP:   "invalid-otp",
	}

	result, err := verifyOTPUC.Execute(context.Background(), input)

	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)

	redis.AssertExpectations(t)
}

func TestVerifyOTP_Integration_EmailMismatch(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("otp2@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "otpuser2", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	otp := "123456"
	redis.On("GetUserFromForgotPasswordOTP", mock.Anything, otp).Return(user.ID, nil).Once()

	verifyOTPUC := usecases.NewVerifyOTPUseCase(userRepo, logger, redis)

	input := usecases.VerifyOTPInput{
		Email: "different@example.com",
		OTP:   otp,
	}

	result, err := verifyOTPUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrInvalidOTP)
	assert.False(t, result.IsValid)

	redis.AssertExpectations(t)
}

