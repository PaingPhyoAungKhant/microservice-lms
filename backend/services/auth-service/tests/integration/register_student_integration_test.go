package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterStudent_Integration_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}
	apiGatewayURL := "http://localhost:3000"

	registerUC := usecases.NewRegisterStudentUseCase(userRepo, publisher, logger, rabbitMQConfig, redis, apiGatewayURL)

	redis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	publisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	input := usecases.RegisterStudentInput{
		Email:    "register-success@example.com",
		Username: "registersuccess",
		Password: "Password123!",
	}

	ctx := context.Background()
	result, err := registerUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Email, result.Email)
	assert.Equal(t, input.Username, result.Username)
	assert.Equal(t, "student", result.Role)
	assert.Equal(t, "active", result.Status)

	createdUser, err := userRepo.FindByEmail(ctx, input.Email)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	assert.Equal(t, input.Email, createdUser.Email.String())
	assert.Equal(t, input.Username, createdUser.Username)
	assert.NoError(t, utils.VerifyPassword(input.Password, createdUser.PasswordHash))

	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestRegisterStudent_Integration_DuplicateEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("existing@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "existing", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	registerUC := usecases.NewRegisterStudentUseCase(userRepo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	input := usecases.RegisterStudentInput{
		Email:    "existing@example.com",
		Username: "newstudent",
		Password: "Password123!",
	}

	_, err = registerUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrEmailAlreadyExists)

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestRegisterStudent_Integration_DuplicateUsername(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("user1@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "existing", role, passwordHash)

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	registerUC := usecases.NewRegisterStudentUseCase(userRepo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	input := usecases.RegisterStudentInput{
		Email:    "user2@example.com",
		Username: "existing",
		Password: "Password123!",
	}

	_, err = registerUC.Execute(ctx, input)

	require.ErrorIs(t, err, usecases.ErrUsernameAlreadyExists)

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestRegisterStudent_Integration_InvalidPassword(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	registerUC := usecases.NewRegisterStudentUseCase(userRepo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	input := usecases.RegisterStudentInput{
		Email:    "register-invalid-password@example.com",
		Username: "registerinvalidpass",
		Password: "weak",
	}

	_, err := registerUC.Execute(context.Background(), input)

	require.Error(t, err)
	// Password validation returns specific error messages
	assert.Contains(t, err.Error(), "password")

	redis.AssertNotCalled(t, "StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything)
	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestRegisterStudent_Integration_MissingRequiredFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := setupUserRepo(t, db)
	publisher := new(mocks.MockPublisher)
	logger := setupTestLogger()
	redis := new(mocks.MockRedis)
	rabbitMQConfig := &config.RabbitMQConfig{}

	registerUC := usecases.NewRegisterStudentUseCase(userRepo, publisher, logger, rabbitMQConfig, redis, "http://localhost:3000")

	tests := []struct {
		name  string
		input usecases.RegisterStudentInput
		err   error
	}{
		{
			name: "missing email",
			input: usecases.RegisterStudentInput{
				Username: "student",
				Password: "Password123!",
			},
			err: usecases.ErrEmailRequired,
		},
		{
			name: "missing username",
			input: usecases.RegisterStudentInput{
				Email:    "student@example.com",
				Password: "Password123!",
			},
			err: usecases.ErrUsernameRequired,
		},
		{
			name: "missing password",
			input: usecases.RegisterStudentInput{
				Email:    "student@example.com",
				Username: "student",
			},
			err: usecases.ErrPasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := registerUC.Execute(context.Background(), tt.input)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

