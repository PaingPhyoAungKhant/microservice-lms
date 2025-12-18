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

func TestLogin_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("login@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "loginuser", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	redis.On("StoreAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	redis.On("StoreRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	redis.On("StoreUserSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	loginUC := usecases.NewLoginUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.LoginInput{
		Email:     "login@example.com",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	}

	result, err := loginUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, user.ID, result.User.ID)
	assert.Equal(t, user.Email.String(), result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestLogin_Integration_InvalidPassword(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("login2@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "loginuser2", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	loginUC := usecases.NewLoginUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.LoginInput{
		Email:     "login2@example.com",
		Password:  "WrongPassword123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	}

	_, err = loginUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrInvalidPassword)

	redis.AssertNotCalled(t, "StoreAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestLogin_Integration_UserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	loginUC := usecases.NewLoginUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.LoginInput{
		Email:     "nonexistent@example.com",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	}

	_, err := loginUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	redis.AssertNotCalled(t, "StoreAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

