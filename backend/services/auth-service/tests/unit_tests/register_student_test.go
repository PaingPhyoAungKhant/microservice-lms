package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
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

func TestRegisterStudent_MissingEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Username: "student",
		Password: "Password123!",
	})
	require.ErrorIs(t, err, usecases.ErrEmailRequired)

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRegisterStudent_MissingUsername(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "student@example.com",
		Password: "Password123!",
	})
	require.ErrorIs(t, err, usecases.ErrUsernameRequired)

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRegisterStudent_MissingPassword(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "student@example.com",
		Username: "student",
	})
	require.ErrorIs(t, err, usecases.ErrPasswordRequired)

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRegisterStudent_InvalidEmail(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "not-an-email",
		Username: "student",
		Password: "Password123!",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")

	repo.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
}

func TestRegisterStudent_EmailAlreadyExists(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("existing@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	existing := entities.NewUser(emailVO, "existing", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	repo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existing, nil).Once()

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "existing@example.com",
		Username: "newstudent",
		Password: "Password123!",
	})
	require.ErrorIs(t, err, usecases.ErrEmailAlreadyExists)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestRegisterStudent_UsernameAlreadyExists(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("existing@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	existing := entities.NewUser(emailVO, "existing", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	repo.On("FindByEmail", mock.Anything, "new@example.com").Return((*entities.User)(nil), nil).Once()
	repo.On("FindByUsername", mock.Anything, "existing").Return(existing, nil).Once()

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "new@example.com",
		Username: "existing",
		Password: "Password123!",
	})
	require.ErrorIs(t, err, usecases.ErrUsernameAlreadyExists)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestRegisterStudent_InvalidPassword(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	repo.On("FindByEmail", mock.Anything, "student@example.com").Return((*entities.User)(nil), nil).Once()
	repo.On("FindByUsername", mock.Anything, "student").Return((*entities.User)(nil), nil).Once()

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	_, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "student@example.com",
		Username: "student",
		Password: "weak",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")

	repo.AssertExpectations(t)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestRegisterStudent_Success(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	logger := logger.NewNop()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}
	apiGatewayURL := "http://localhost:3000"

	var createdUser *entities.User
	repo.On("FindByEmail", mock.Anything, "student@example.com").Return((*entities.User)(nil), nil).Once()
	repo.On("FindByUsername", mock.Anything, "student").Return((*entities.User)(nil), nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
		createdUser = u
		return true
	})).Return(nil).Once()
	redis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeAuthStudentRegistered, mock.Anything).Return(nil).Once()

	uc := usecases.NewRegisterStudentUseCase(repo, publisher, logger, rabbitMQConfig, redis, apiGatewayURL)

	result, err := uc.Execute(context.Background(), usecases.RegisterStudentInput{
		Email:    "student@example.com",
		Username: "student",
		Password: "Password123!",
	})
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, result)
	assert.Equal(t, "student@example.com", createdUser.Email.String())
	assert.Equal(t, "student", createdUser.Username)
	assert.Equal(t, "student", createdUser.Role.String())
	assert.NoError(t, utils.VerifyPassword("Password123!", createdUser.PasswordHash))

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.AuthStudentRegisteredEvent)
	require.True(t, ok)
	assert.Equal(t, createdUser.ID, event.ID)
	assert.Equal(t, "student@example.com", event.Email)
	assert.Equal(t, "student", event.Username)
	assert.Contains(t, event.EmailVerificationURL, apiGatewayURL)

	repo.AssertExpectations(t)
	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

