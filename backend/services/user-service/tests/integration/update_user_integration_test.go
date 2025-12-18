package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser_Integration_Success(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	updateUserUC := usecases.NewUpdateUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("update@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "updateuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	newEmail := "updated@example.com"
	newUsername := "updateduser"
	newRole := "instructor"
	newStatus := "active"

	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	input := usecases.UpdateUserInput{
		ID:       user.ID,
		Email:    &newEmail,
		Username: &newUsername,
		Role:     &newRole,
		Status:   &newStatus,
	}

	result, err := updateUserUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, newEmail, result.Email)
	assert.Equal(t, newUsername, result.Username)
	assert.Equal(t, newRole, result.Role)
	assert.Equal(t, newStatus, result.Status)

	updatedUser, err := userRepo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, newEmail, updatedUser.Email.String())
	assert.Equal(t, newUsername, updatedUser.Username)

	mockPublisher.AssertExpectations(t)
}

func TestUpdateUser_Integration_PartialUpdate(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	updateUserUC := usecases.NewUpdateUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("partial@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "partialuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	newEmail := "partialupdated@example.com"
	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	input := usecases.UpdateUserInput{
		ID:    user.ID,
		Email: &newEmail,
	}

	result, err := updateUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, newEmail, result.Email)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Role.String(), result.Role)

	mockPublisher.AssertExpectations(t)
}

func TestUpdateUser_Integration_InvalidEmail(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	updateUserUC := usecases.NewUpdateUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("invalid@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "invaliduser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	invalidEmail := "not-an-email"
	input := usecases.UpdateUserInput{
		ID:    user.ID,
		Email: &invalidEmail,
	}

	_, err = updateUserUC.Execute(ctx, input)

	require.Error(t, err)
}

func TestUpdateUser_Integration_InvalidRole(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	updateUserUC := usecases.NewUpdateUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("invalidrole@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "invalidroleuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	invalidRole := "invalid-role"
	input := usecases.UpdateUserInput{
		ID:   user.ID,
		Role: &invalidRole,
	}

	_, err = updateUserUC.Execute(ctx, input)

	require.Error(t, err)
}

