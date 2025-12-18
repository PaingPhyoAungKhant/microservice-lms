package unit_test

import (
	"context"
	"errors"
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

func TestForgotPassword_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewForgotPasswordUseCase(repo, publisher, logger)

	_, err := uc.Execute(context.Background(), usecases.ForgotPasswordInput{})
	require.ErrorIs(t, err, usecases.ErrEmailRequired)

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestForgotPassword_InvalidEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	uc := usecases.NewForgotPasswordUseCase(repo, publisher, logger)

	_, err := uc.Execute(context.Background(), usecases.ForgotPasswordInput{
		Email: "not-an-email",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewForgotPasswordUseCase(repo, publisher, logger)

	_, err := uc.Execute(context.Background(), usecases.ForgotPasswordInput{
		Email: "user@example.com",
	})
	require.ErrorIs(t, err, usecases.ErrUserNotFound)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestForgotPassword_PublisherFailure(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserForgotPassword, mock.Anything).Return(errors.New("publisher error")).Once()

	uc := usecases.NewForgotPasswordUseCase(repo, publisher, logger)

	_, err := uc.Execute(context.Background(), usecases.ForgotPasswordInput{
		Email: "user@example.com",
	})
	require.ErrorIs(t, err, usecases.ErrInternalServerError)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestForgotPassword_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserForgotPassword, mock.Anything).Return(nil).Once()

	uc := usecases.NewForgotPasswordUseCase(repo, publisher, logger)

	result, err := uc.Execute(context.Background(), usecases.ForgotPasswordInput{
		Email: "user@example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Password reset email sent", result.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthUserForgotPasswordEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)
	assert.Equal(t, "user@example.com", event.Email)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

