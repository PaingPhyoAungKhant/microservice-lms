package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
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

func TestCreateUser_InvalidEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	uc := usecases.NewCreateUserUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.CreateUserInput{
		Email:    "not-an-email",
		Username: "tester",
		Password: "Secret123!",
		Role:     "student",
	})
	require.ErrorIs(t, err, valueobjects.ErrInvalidEmail)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateUser_InvalidRole(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	uc := usecases.NewCreateUserUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.CreateUserInput{
		Email:    "user@example.com",
		Username: "tester",
		Password: "Secret123!",
		Role:     "invalid-role",
	})
	require.ErrorIs(t, err, valueobjects.ErrInvalidRole)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("existing@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	existing := entities.NewUser(emailVO, "existing", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existing, nil).Once()
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	uc := usecases.NewCreateUserUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.CreateUserInput{
		Email:    "existing@example.com",
		Username: "someone",
		Password: "Secret123!",
		Role:     "student",
	})
	require.ErrorIs(t, err, usecases.ErrEmailAlreadyExists)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("existing@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	existing := entities.NewUser(emailVO, "existing", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "new@example.com").Return((*entities.User)(nil), nil).Once()
	repo.On("FindByUsername", mock.Anything, "existing").Return(existing, nil).Once()
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	uc := usecases.NewCreateUserUseCase(repo, publisher, logger, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.CreateUserInput{
		Email:    "new@example.com",
		Username: "existing",
		Password: "Secret123!",
		Role:     "student",
	})
	require.ErrorIs(t, err, usecases.ErrUsernameAlreadyExists)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	apiGatewayURL := "http://localhost:3000"
	input := usecases.CreateUserInput{
		Email:    "newuser@example.com",
		Username: "newuser",
		Password: "Secret123!",
		Role:     "student",
	}

	var createdUser *entities.User
	repo.On("FindByEmail", mock.Anything, input.Email).Return((*entities.User)(nil), nil).Once()
	repo.On("FindByUsername", mock.Anything, input.Username).Return((*entities.User)(nil), nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
		createdUser = u
		return true
	})).Return(nil).Once()
	redis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeUserCreated, mock.Anything).Return(nil).Once()

	uc := usecases.NewCreateUserUseCase(repo, publisher, logger, redis, apiGatewayURL)

	dto, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Equal(t, input.Email, createdUser.Email.String())
	require.Equal(t, input.Username, createdUser.Username)
	require.NoError(t, utils.VerifyPassword(input.Password, createdUser.PasswordHash))

	require.Len(t, publisher.Calls, 1)
	call := publisher.Calls[0]
	event, ok := call.Arguments.Get(2).(events.UserCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, createdUser.ID, event.ID)
	assert.Equal(t, input.Email, event.Email)
	assert.Equal(t, input.Username, event.Username)
	assert.Equal(t, createdUser.EmailVerified, event.EmailVerified)
	assert.Equal(t, createdUser.EmailVerifiedAt, event.EmailVerifiedAt)
	assert.Contains(t, event.EmailVerificationURL, apiGatewayURL)
	assert.Contains(t, event.EmailVerificationURL, "/api/v1/auth/verify-email?token=")

	assertDTOEqualUser(t, dto, createdUser)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

