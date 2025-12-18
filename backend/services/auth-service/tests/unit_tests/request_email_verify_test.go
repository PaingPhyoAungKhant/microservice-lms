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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequestEmailVerify_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewRequestEmailVerifyUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RequestEmailVerifyInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRequestEmailVerify_InvalidEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	uc := usecases.NewRequestEmailVerifyUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RequestEmailVerifyInput{
		Email: "not-an-email",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRequestEmailVerify_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewRequestEmailVerifyUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RequestEmailVerifyInput{
		Email: "user@example.com",
	})
	require.ErrorIs(t, err, usecases.ErrUserNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestRequestEmailVerify_AlreadyVerified(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.EmailVerified = true

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()

	uc := usecases.NewRequestEmailVerifyUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	result, err := uc.Execute(context.Background(), usecases.RequestEmailVerifyInput{
		Email: "user@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "Email already verified", result.Message)

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestRequestEmailVerify_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.EmailVerified = false

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	apiGatewayURL := "http://localhost:3000"

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()
	redis.On("StoreVerifyEmailToken", mock.Anything, user.ID, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserRequestedEmailVerification, mock.Anything).Return(nil).Once()

	uc := usecases.NewRequestEmailVerifyUseCase(repo, publisher, logger, redis, apiGatewayURL)

	result, err := uc.Execute(context.Background(), usecases.RequestEmailVerifyInput{
		Email: "user@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "Verification email sent", result.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthUserRequestedEmailVerificationEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)
	assert.Equal(t, "user@example.com", event.Email)
	assert.Contains(t, event.EmailVerificationURL, apiGatewayURL)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

