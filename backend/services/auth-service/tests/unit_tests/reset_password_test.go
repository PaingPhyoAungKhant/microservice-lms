package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResetPassword_MissingToken(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		NewPassword: "NewPassword123!",
		IPAddress:   "127.0.0.1",
		UserAgent:   "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidPasswordResetToken)
	assert.Contains(t, result.Message, "invalid password reset token")

	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_MissingPassword(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:     "valid-token",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidPassword)
	assert.Contains(t, result.Message, "invalid password")

	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_MissingIPAddress(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "valid-token",
		NewPassword: "NewPassword123!",
		UserAgent:  "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidIPAddress)
	assert.Contains(t, result.Message, "invalid IP address")

	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_MissingUserAgent(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "valid-token",
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidUserAgent)
	assert.Contains(t, result.Message, "invalid user agent")

	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_InvalidToken(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromResetPasswordToken", mock.Anything, "invalid-token").Return("", nil).Once()

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "invalid-token",
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidPasswordResetToken)
	assert.Contains(t, result.Message, "invalid password reset token")

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
}

func TestResetPassword_InvalidNewPassword(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "valid-token",
		NewPassword: "weak",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	})
	require.Error(t, err)
	// Password validation returns specific error messages, not generic "invalid password"
	assert.Contains(t, result.Message, "password")

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	// Password validation happens before redis call, so redis shouldn't be called
	redis.AssertNotCalled(t, "GetUserFromResetPasswordToken", mock.Anything, mock.Anything)
}

func TestResetPassword_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromResetPasswordToken", mock.Anything, "valid-token").Return("user-id", nil).Once()
	redis.On("RevokeResetPasswordToken", mock.Anything, "valid-token").Return(nil).Once()
	repo.On("FindByID", mock.Anything, "user-id").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "valid-token",
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrUserNotFound)
	assert.Contains(t, result.Message, "user not found")

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	// Use a valid bcrypt hash for the old password
	oldHash, _ := utils.HashPassword("OldPassword123!")
	user := entities.NewUser(emailVO, "user", roleVO, oldHash)

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromResetPasswordToken", mock.Anything, "valid-token").Return(user.ID, nil).Once()
	redis.On("RevokeResetPasswordToken", mock.Anything, "valid-token").Return(nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	repo.On("UpdatePassword", mock.Anything, user.ID, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserResetPassword, mock.Anything).Return(nil).Once()

	uc := usecases.NewResetPasswordUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.ResetPasswordInput{
		Token:      "valid-token",
		NewPassword: "NewPassword123!",
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Password reset successful", result.Message)

	var updatedPasswordHash string
	for _, call := range repo.Calls {
		if call.Method == "UpdatePassword" {
			// UpdatePassword signature: (ctx, userID, passwordHash)
			// So arguments are: [0]=ctx, [1]=userID, [2]=passwordHash
			updatedPasswordHash = call.Arguments.Get(2).(string)
			break
		}
	}
	assert.NoError(t, utils.VerifyPassword("NewPassword123!", updatedPasswordHash))

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthUserResetPasswordEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

