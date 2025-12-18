package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Integration_Success(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	apiGatewayURL := "http://localhost:3000"
	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		apiGatewayURL,
	)

	mockRedis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	input := usecases.CreateUserInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "Password123!",
		Role:     "student",
	}

	ctx := context.Background()
	result, err := createUserUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, input.Email, result.Email)
	assert.Equal(t, input.Username, result.Username)
	assert.Equal(t, input.Role, result.Role)
	assert.Equal(t, "active", result.Status)
	assert.False(t, result.EmailVerified)

	createdUser, err := userRepo.FindByEmail(ctx, input.Email)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	assert.Equal(t, input.Email, createdUser.Email.String())
	assert.Equal(t, input.Username, createdUser.Username)

	mockRedis.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCreateUser_Integration_DuplicateEmail(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		"http://localhost:3000",
	)

	ctx := context.Background()

	input1 := usecases.CreateUserInput{
		Email:    "user@example.com",
		Username: "user1",
		Password: "Password123!",
		Role:     "student",
	}

	mockRedis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err = createUserUC.Execute(ctx, input1)
	require.NoError(t, err)

	input2 := usecases.CreateUserInput{
		Email:    "user@example.com",
		Username: "user2",
		Password: "Password123!",
		Role:     "student",
	}

	_, err = createUserUC.Execute(ctx, input2)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrEmailAlreadyExists, err)
}

func TestCreateUser_Integration_DuplicateUsername(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		"http://localhost:3000",
	)

	ctx := context.Background()

	input1 := usecases.CreateUserInput{
		Email:    "user1@example.com",
		Username: "duplicateuser",
		Password: "Password123!",
		Role:     "student",
	}

	mockRedis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err = createUserUC.Execute(ctx, input1)
	require.NoError(t, err)

	input2 := usecases.CreateUserInput{
		Email:    "user2@example.com",
		Username: "duplicateuser",
		Password: "Password123!",
		Role:     "student",
	}

	_, err = createUserUC.Execute(ctx, input2)
	require.Error(t, err)
	assert.Equal(t, usecases.ErrUsernameAlreadyExists, err)
}

func TestCreateUser_Integration_InvalidPassword(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		"http://localhost:3000",
	)

	input := usecases.CreateUserInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "asd",
		Role:     "student",
	}

	ctx := context.Background()
	_, err = createUserUC.Execute(ctx, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestCreateUser_Integration_AllRoles(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		"http://localhost:3000",
	)

	ctx := context.Background()
	roles := []string{"student", "instructor", "admin"}

	for i, role := range roles {
		mockRedis.On("StoreVerifyEmailToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		input := usecases.CreateUserInput{
			Email:    fmt.Sprintf("user%d@example.com", i),
			Username: fmt.Sprintf("user%d", i),
			Password: "Password123!",
			Role:     role,
		}

		result, err := createUserUC.Execute(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, role, result.Role)
	}
}

func TestCreateUser_Integration_MissingRequiredFields(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	mockRedis := new(mocks.MockRedis)
	logger := logger.NewNop()

	createUserUC := usecases.NewCreateUserUseCase(
		userRepo,
		mockPublisher,
		logger,
		mockRedis,
		"http://localhost:3000",
	)

	ctx := context.Background()

	tests := []struct {
		name  string
		input usecases.CreateUserInput
		err   error
	}{
		{
			name: "missing email",
			input: usecases.CreateUserInput{
				Username: "testuser",
				Password: "Password123!",
				Role:     "student",
			},
			err: usecases.ErrEmailRequired,
		},
		{
			name: "missing username",
			input: usecases.CreateUserInput{
				Email:    "test@example.com",
				Password: "Password123!",
				Role:     "student",
			},
			err: usecases.ErrUsernameRequired,
		},
		{
			name: "missing password",
			input: usecases.CreateUserInput{
				Email:    "test@example.com",
				Username: "testuser",
				Role:     "student",
			},
			err: usecases.ErrPasswordRequired,
		},
		{
			name: "missing role",
			input: usecases.CreateUserInput{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "Password123!",
			},
			err: usecases.ErrRoleRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createUserUC.Execute(ctx, tt.input)
			require.Error(t, err)
			assert.Equal(t, tt.err, err)
		})
	}
}

