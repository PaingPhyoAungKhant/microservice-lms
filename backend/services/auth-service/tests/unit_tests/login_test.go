package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func setupJwtManagerForUnit() *utils.JwtManager {
	return utils.NewJwtManager("test-secret", 15*time.Minute, 24*time.Hour)
}

func TestLogin_InvalidEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "not-an-email",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestLogin_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.LoginInput{
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrEmailRequired)
}

func TestLogin_MissingPassword(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "user@example.com",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrPasswordRequired)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "user@example.com",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrUserNotFound)

	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("CorrectPassword123!")
	user := entities.NewUser(emailVO, "user", roleVO, passwordHash)

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	_, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "user@example.com",
		Password:  "WrongPassword123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.ErrorIs(t, err, usecases.ErrInvalidPassword)

	repo.AssertExpectations(t)
}


func TestLogin_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(emailVO, "user", roleVO, passwordHash)

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()
	redis.On("StoreAccessToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil).Once()
	redis.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil).Once()
	redis.On("StoreUserSession", mock.Anything, mock.Anything, user.ID, mock.Anything, mock.Anything, "127.0.0.1", "test-agent", mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserLoggedIn, mock.Anything).Return(nil).Once()

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	result, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "user@example.com",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, user.ID, result.User.ID)
	assert.Equal(t, user.Email.String(), result.User.Email)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestLogin_RedisStorageFailure(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(emailVO, "user", roleVO, passwordHash)

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	jwtManager := setupJwtManagerForUnit()
	redis := new(mocks.MockRedis)

	repo.On("FindByEmail", mock.Anything, "user@example.com").Return(user, nil).Once()
	redis.On("StoreAccessToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(errors.New("redis error")).Once()
	redis.On("StoreRefreshToken", mock.Anything, user.ID, mock.Anything, mock.Anything).Return(errors.New("redis error")).Once()
	redis.On("StoreUserSession", mock.Anything, mock.Anything, user.ID, mock.Anything, mock.Anything, "127.0.0.1", "test-agent", mock.Anything).Return(errors.New("redis error")).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthUserLoggedIn, mock.Anything).Return(nil).Once()

	uc := usecases.NewLoginUseCase(repo, publisher, logger, jwtManager, redis)

	result, err := uc.Execute(context.Background(), usecases.LoginInput{
		Email:     "user@example.com",
		Password:  "Password123!",
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

