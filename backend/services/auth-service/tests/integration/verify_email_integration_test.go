package integration

import (
	"context"
	"errors"
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

func TestVerifyEmail_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("verifyemail@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "verifyemailuser", role, passwordHash)
	user.EmailVerified = false

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "valid-verify-token"
	redis.On("GetUserFromVerifyEmailToken", mock.Anything, token).Return(user.ID, nil).Once()
	redis.On("RevokeVerifyEmailToken", mock.Anything, token).Return(nil).Once()

	verifyEmailUC := usecases.NewVerifyEmailUseCase(userRepo, publisher, logger, redis)

	input := usecases.VerifyEmailInput{
		Email: "verifyemail@example.com",
		Token: token,
	}

	result, err := verifyEmailUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Email verified successfully", result.Message)

	updatedUser, err := userRepo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, updatedUser.EmailVerified)

	redis.AssertExpectations(t)
}

func TestVerifyEmail_Integration_InvalidToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	// When Redis returns an error, the use case treats it as invalid token
	redis.On("GetUserFromVerifyEmailToken", mock.Anything, "invalid-token").Return("", errors.New("token not found")).Once()

	verifyEmailUC := usecases.NewVerifyEmailUseCase(userRepo, publisher, logger, redis)

	input := usecases.VerifyEmailInput{
		Email: "user@example.com",
		Token: "invalid-token",
	}

	_, err := verifyEmailUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")

	redis.AssertExpectations(t)
}

func TestVerifyEmail_Integration_AlreadyVerified(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("verified@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "verifieduser", role, passwordHash)
	user.EmailVerified = true

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "valid-token"
	redis.On("GetUserFromVerifyEmailToken", mock.Anything, token).Return(user.ID, nil).Once()

	verifyEmailUC := usecases.NewVerifyEmailUseCase(userRepo, publisher, logger, redis)

	input := usecases.VerifyEmailInput{
		Email: "verified@example.com",
		Token: token,
	}

	result, err := verifyEmailUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Email already verified", result.Message)

	redis.AssertExpectations(t)
}

func TestVerifyEmail_Integration_EmailMismatch(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("user@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "user", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "valid-token"
	redis.On("GetUserFromVerifyEmailToken", mock.Anything, token).Return(user.ID, nil).Once()

	verifyEmailUC := usecases.NewVerifyEmailUseCase(userRepo, publisher, logger, redis)

	input := usecases.VerifyEmailInput{
		Email: "different@example.com",
		Token: token,
	}

	_, err = verifyEmailUC.Execute(ctx, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "email does not match")

	redis.AssertExpectations(t)
}

