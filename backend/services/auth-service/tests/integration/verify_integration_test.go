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

func TestVerify_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("verify@example.com")
	role, _ := valueobjects.NewRole("student")
	activeStatus, _ := valueobjects.NewStatus("active")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "verifyuser", role, passwordHash)
	user.Status = activeStatus

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	require.NoError(t, err)

	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return(user.ID, nil).Once()

	verifyUC := usecases.NewVerifyUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	}

	result, err := verifyUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email.String(), result.Email)
	assert.Equal(t, user.Role.String(), result.Role)

	redis.AssertExpectations(t)
}

func TestVerify_Integration_InsufficientPermissions(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("verify2@example.com")
	role, _ := valueobjects.NewRole("student")
	activeStatus, _ := valueobjects.NewStatus("active")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "verifyuser2", role, passwordHash)
	user.Status = activeStatus

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	require.NoError(t, err)

	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return(user.ID, nil).Once()

	verifyUC := usecases.NewVerifyUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "admin",
	}

	_, err = verifyUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrInsufficientPermissions)

	redis.AssertExpectations(t)
}

func TestVerify_Integration_InvalidToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	verifyUC := usecases.NewVerifyUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.VerifyInput{
		Token:        "invalid-token",
		RequiredRole: "student",
	}

	_, err := verifyUC.Execute(context.Background(), input)

	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	redis.AssertNotCalled(t, "GetUserFromAccessToken", mock.Anything, mock.Anything)
}

func TestVerify_Integration_InactiveUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	jwtManager := setupJwtManager()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("verify3@example.com")
	role, _ := valueobjects.NewRole("student")
	pendingStatus, _ := valueobjects.NewStatus("pending")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "verifyuser3", role, passwordHash)
	user.Status = pendingStatus

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	require.NoError(t, err)

	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return(user.ID, nil).Once()

	verifyUC := usecases.NewVerifyUseCase(userRepo, publisher, logger, jwtManager, redis)

	input := usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	}

	_, err = verifyUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	redis.AssertExpectations(t)
}

