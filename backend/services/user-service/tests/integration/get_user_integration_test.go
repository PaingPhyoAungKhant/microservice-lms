package integration

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser_Integration_Success(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	getUserUC := usecases.NewGetUserUseCase(userRepo, logger)

	ctx := context.Background()

	email, _ := valueobjects.NewEmail("getuser@example.com")
	role, _ := valueobjects.NewRole("student")
	passwordHash, _ := utils.HashPassword("Password123!")
	user := entities.NewUser(email, "getuser", role, passwordHash)

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	input := usecases.GetUserInput{
		UserID: user.ID,
	}

	result, err := getUserUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email.String(), result.Email)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Role.String(), result.Role)
}

func TestGetUser_Integration_NotFound(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	getUserUC := usecases.NewGetUserUseCase(userRepo, logger)

	ctx := context.Background()

	nonExistentUUID := "00000000-0000-0000-0000-000000000000"
	input := usecases.GetUserInput{
		UserID: nonExistentUUID,
	}

	_, err = getUserUC.Execute(ctx, input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrUserNotFound, err)
}

