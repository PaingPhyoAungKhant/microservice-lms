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

func TestForgotPassword_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("forgot@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "forgotuser", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserForgotPassword, mock.Anything).Return(nil).Once()

	forgotPasswordUC := usecases.NewForgotPasswordUseCase(userRepo, publisher, logger)

	input := usecases.ForgotPasswordInput{
		Email: "forgot@example.com",
	}

	result, err := forgotPasswordUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Password reset email sent", result.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthUserForgotPasswordEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)
	assert.Equal(t, "forgot@example.com", event.Email)

	publisher.AssertExpectations(t)
}

func TestForgotPassword_Integration_UserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()

	forgotPasswordUC := usecases.NewForgotPasswordUseCase(userRepo, publisher, logger)

	input := usecases.ForgotPasswordInput{
		Email: "nonexistent@example.com",
	}

	_, err := forgotPasswordUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestForgotPassword_Integration_InvalidEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()

	forgotPasswordUC := usecases.NewForgotPasswordUseCase(userRepo, publisher, logger)

	input := usecases.ForgotPasswordInput{
		Email: "not-an-email",
	}

	_, err := forgotPasswordUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

