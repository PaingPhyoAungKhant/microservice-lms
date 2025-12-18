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

func TestDeleteUser_Integration_Success(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	deleteUserUC := usecases.NewDeleteUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("delete@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "deleteuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	input := usecases.DeleteUserInput{
		UserID: user.ID,
	}

	result, err := deleteUserUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Message, "Deleted")

	_, err = userRepo.FindByID(ctx, user.ID)
	require.Error(t, err)

	mockPublisher.AssertExpectations(t)
}

func TestDeleteUser_Integration_CannotDeleteAdmin(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	deleteUserUC := usecases.NewDeleteUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("admin@example.com")
	role, _ := valueobjects.NewRole("admin")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "adminuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	input := usecases.DeleteUserInput{
		UserID: user.ID,
	}

	_, err = deleteUserUC.Execute(ctx, input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrCannotDeleteAdminUser, err)

	existingUser, err := userRepo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, existingUser)
}

func TestDeleteUser_Integration_NotFound(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	mockPublisher := new(mocks.MockPublisher)
	logger := logger.NewNop()

	deleteUserUC := usecases.NewDeleteUserUseCase(userRepo, mockPublisher, logger)

	ctx := context.Background()

	input := usecases.DeleteUserInput{
		UserID: "non-existent-id",
	}

	_, err = deleteUserUC.Execute(ctx, input)

	require.Error(t, err)
}

