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

func TestVerify_InvalidToken(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        "invalid-token",
		RequiredRole: "student",
	})
	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	redis.AssertNotCalled(t, "GetUserFromAccessToken", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestVerify_TokenNotFoundInRedis(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	accessToken, _ := jwtManager.GenerateAccessToken("user-id", "user@example.com", "student", "active")
	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return("", errors.New("not found")).Once()

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	})
	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	redis.AssertExpectations(t)
}

func TestVerify_InsufficientPermissions_RoleMismatch(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	statusVO, _ := valueobjects.NewStatus("active")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.Status = statusVO

	accessToken, _ := jwtManager.GenerateAccessToken("user-id", "user@example.com", "student", "active")
	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return("user-id", nil).Once()
	repo.On("FindByID", mock.Anything, "user-id").Return(user, nil).Once()

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "admin",
	})
	require.ErrorIs(t, err, usecases.ErrInsufficientPermissions)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
}

func TestVerify_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	accessToken, _ := jwtManager.GenerateAccessToken("user-id", "user@example.com", "student", "active")
	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return("user-id", nil).Once()
	repo.On("FindByID", mock.Anything, "user-id").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	})
	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerify_InactiveUser(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	pendingStatus, _ := valueobjects.NewStatus("pending")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.Status = pendingStatus

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	accessToken, _ := jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	})
	require.ErrorIs(t, err, usecases.ErrUnauthorized)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestVerify_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	activeStatus, _ := valueobjects.NewStatus("active")
	user := entities.NewUser(emailVO, "user", roleVO, "hash")
	user.Status = activeStatus

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	accessToken, _ := jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	redis.On("GetUserFromAccessToken", mock.Anything, accessToken).Return(user.ID, nil).Once()
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewVerifyUseCase(repo, publisher, logger, jwtManager, redis)

	result, err := uc.Execute(context.Background(), usecases.VerifyInput{
		Token:        accessToken,
		RequiredRole: "student",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email.String(), result.Email)
	assert.Equal(t, user.Role.String(), result.Role)

	redis.AssertExpectations(t)
	repo.AssertExpectations(t)
}

