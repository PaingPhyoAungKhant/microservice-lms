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

func TestRequestEmailVerify_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	apiGatewayURL := "http://localhost:3000"

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("requestverify@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "requestverifyuser", role, passwordHash)
	user.EmailVerified = false

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	redis.On("StoreVerifyEmailToken", mock.Anything, user.ID, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserRequestedEmailVerification, mock.Anything).Return(nil).Once()

	requestEmailVerifyUC := usecases.NewRequestEmailVerifyUseCase(userRepo, publisher, logger, redis, apiGatewayURL)

	input := usecases.RequestEmailVerifyInput{
		Email: "requestverify@example.com",
	}

	result, err := requestEmailVerifyUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Verification email sent", result.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthUserRequestedEmailVerificationEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)
	assert.Equal(t, "requestverify@example.com", event.Email)
	assert.Contains(t, event.EmailVerificationURL, apiGatewayURL)

	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestRequestEmailVerify_Integration_UserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	requestEmailVerifyUC := usecases.NewRequestEmailVerifyUseCase(userRepo, publisher, logger, redis, "http://localhost:3000")

	input := usecases.RequestEmailVerifyInput{
		Email: "nonexistent@example.com",
	}

	_, err := requestEmailVerifyUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestRequestEmailVerify_Integration_AlreadyVerified(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("alreadyverified@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "alreadyverifieduser", role, passwordHash)
	user.EmailVerified = true

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	requestEmailVerifyUC := usecases.NewRequestEmailVerifyUseCase(userRepo, publisher, logger, redis, "http://localhost:3000")

	input := usecases.RequestEmailVerifyInput{
		Email: "alreadyverified@example.com",
	}

	result, err := requestEmailVerifyUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Email already verified", result.Message)

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

