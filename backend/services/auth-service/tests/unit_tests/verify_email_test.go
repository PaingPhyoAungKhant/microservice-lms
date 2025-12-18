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

func TestVerifyEmail_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	// The use case doesn't validate email presence - it's optional
	// The token lookup will happen, then email validation if provided
	// So we need to mock the redis call
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "valid-token").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	repo.On("UpdateEmailVerified", mock.Anything, user.ID, true).Return(nil).Once()
	redis.On("RevokeVerifyEmailToken", mock.Anything, "valid-token").Return(nil).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Token: "valid-token",
	})
	// Email is optional, so verification should succeed
	require.NoError(t, err)
	assert.Equal(t, "Email verified successfully", result.Message)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyEmail_MissingToken(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "user@example.com",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token is required")

	redis.AssertNotCalled(t, "GetUserFromVerifyEmailToken", mock.Anything, mock.Anything)
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "invalid-token").Return("", errors.New("not found")).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "user@example.com",
		Token: "invalid-token",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
}

func TestVerifyEmail_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "valid-token").Return("user-id", nil).Once()
	repo.On("FindByID", mock.Anything, "user-id").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "user@example.com",
		Token: "valid-token",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyEmail_EmailMismatch(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "valid-token").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "different@example.com",
		Token: "valid-token",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email does not match")

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyEmail_AlreadyVerified(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.EmailVerified = true

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "valid-token").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "user@example.com",
		Token: "valid-token",
	})
	require.NoError(t, err)
	assert.Equal(t, "Email already verified", result.Message)

	repo.AssertNotCalled(t, "UpdateEmailVerified", mock.Anything, mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerifyEmail_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.EmailVerified = false

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "valid-token").Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	repo.On("UpdateEmailVerified", mock.Anything, user.ID, true).Return(nil).Once()
	redis.On("RevokeVerifyEmailToken", mock.Anything, "valid-token").Return(nil).Once()

	uc := usecases.NewVerifyEmailUseCase(repo, publisher, logger, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyEmailInput{
		Email: "user@example.com",
		Token: "valid-token",
	})
	require.NoError(t, err)
	assert.Equal(t, "Email verified successfully", result.Message)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

