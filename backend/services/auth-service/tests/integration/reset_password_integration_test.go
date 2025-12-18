package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResetPassword_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("reset@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("OldPassword123!")
	user := entities.NewUser(email, "resetuser", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "valid-reset-token"
	redis.On("GetUserFromResetPasswordToken", mock.Anything, token).Return(user.ID, nil).Once()
	redis.On("RevokeResetPasswordToken", mock.Anything, token).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserResetPassword, mock.Anything).Return(nil).Once()

	resetPasswordUC := usecases.NewResetPasswordUseCase(userRepo, publisher, logger, redis)

	input := usecases.ResetPasswordInput{
		Token:      token,
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	}

	result, err := resetPasswordUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Password reset successful", result.Message)

	updatedUser, err := userRepo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.NoError(t, utils.VerifyPassword("NewPassword123!", updatedUser.PasswordHash))

	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestResetPassword_Integration_InvalidToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromResetPasswordToken", mock.Anything, "invalid-token").Return("", nil).Once()

	resetPasswordUC := usecases.NewResetPasswordUseCase(userRepo, publisher, logger, redis)

	input := usecases.ResetPasswordInput{
		Token:      "invalid-token",
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	}

	result, err := resetPasswordUC.Execute(context.Background(), input)

	require.ErrorIs(t, err, usecases.ErrInvalidPasswordResetToken)
	assert.Contains(t, result.Message, "invalid password reset token")

	redis.AssertExpectations(t)
}

func TestResetPassword_Integration_InvalidPassword(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("reset2@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("OldPassword123!")
	user := entities.NewUser(email, "resetuser2", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "valid-token"
	redis.On("GetUserFromResetPasswordToken", mock.Anything, token).Return(user.ID, nil).Once()

	resetPasswordUC := usecases.NewResetPasswordUseCase(userRepo, publisher, logger, redis)

	input := usecases.ResetPasswordInput{
		Token:      token,
		NewPassword: "weak",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	}

	result, err := resetPasswordUC.Execute(ctx, input)

	require.Error(t, err)
	// Password validation returns specific error messages, not generic "invalid password"
	assert.Contains(t, result.Message, "password")

	// Password validation happens before redis call, so redis shouldn't be called
	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_Integration_UserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	nonExistentUserID := "00000000-0000-0000-0000-000000000000"
	token := "valid-token"
	redis.On("GetUserFromResetPasswordToken", mock.Anything, token).Return(nonExistentUserID, nil).Once()
	redis.On("RevokeResetPasswordToken", mock.Anything, token).Return(nil).Once()

	resetPasswordUC := usecases.NewResetPasswordUseCase(userRepo, publisher, logger, redis)

	input := usecases.ResetPasswordInput{
		Token:      token,
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	}

	result, err := resetPasswordUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, result.Message, "user not found")

	redis.AssertExpectations(t)
}

