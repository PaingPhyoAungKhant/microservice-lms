package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindUser_Integration_ByRole(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	var _ *sql.DB = db
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@asto-lms.local'")
	require.NoError(t, err)

	roles := []string{"student", "instructor", "admin"}
	for _, roleStr := range roles {
		email, _ := valueobjects.NewEmail(fmt.Sprintf("find%s@example.com", roleStr))
		role, _ := valueobjects.NewRole(roleStr)
		passwordHash, _ := utils.HashPassword("Password123!")
		user := entities.NewUser(email, fmt.Sprintf("find%s", roleStr), role, passwordHash)
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	studentRole := valueobjects.Role("student")
	limit := 10
	input := usecases.FindUserInput{
		Role:  &studentRole,
		Limit: &limit,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "student", result.Users[0].Role)
}

func TestFindUser_Integration_ByStatus(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	var _ *sql.DB = db
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@asto-lms.local'")
	require.NoError(t, err)

	statuses := []string{"pending", "active", "inactive"}
	for i, statusStr := range statuses {
		email, _ := valueobjects.NewEmail(fmt.Sprintf("status%d@example.com", i))
		role, _ := valueobjects.NewRole("student")
		passwordHash, _ := utils.HashPassword("Password123!")
		user := entities.NewUser(email, fmt.Sprintf("status%d", i), role, passwordHash)
		status, _ := valueobjects.NewStatus(statusStr)
		user.Status = status
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	activeStatus := valueobjects.Status("active")
	limit := 10
	input := usecases.FindUserInput{
		Status: &activeStatus,
		Limit:  &limit,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "active", result.Users[0].Status)
}

func TestFindUser_Integration_WithSearch(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	emails := []string{"search1@example.com", "search2@example.com", "other@example.com"}
	for _, emailStr := range emails {
		email, _ := valueobjects.NewEmail(emailStr)
		role, _ := valueobjects.NewRole("student")
		passwordHash, _ := utils.HashPassword("Password123!")
		user := entities.NewUser(email, "user"+emailStr, role, passwordHash)
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	searchQuery := "search"
	limit := 10
	input := usecases.FindUserInput{
		SearchQuery: &searchQuery,
		Limit:       &limit,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Users, 2)
}

func TestFindUser_Integration_WithPagination(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	var _ *sql.DB = db
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@asto-lms.local'")
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		email, _ := valueobjects.NewEmail(fmt.Sprintf("paginate%d@example.com", i))
		role, _ := valueobjects.NewRole("student")
		passwordHash, _ := utils.HashPassword("Password123!")
		user := entities.NewUser(email, fmt.Sprintf("paginate%d", i), role, passwordHash)
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	limit := 2
	offset := 0
	input := usecases.FindUserInput{
		Limit:  &limit,
		Offset: &offset,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 5, result.Total)
	assert.Len(t, result.Users, 2)

	offset = 2
	input.Offset = &offset

	result2, err := findUserUC.Execute(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, 5, result2.Total)
	assert.Len(t, result2.Users, 2)
}

func TestFindUser_Integration_WithSorting(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	var _ *sql.DB = db
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@asto-lms.local'")
	require.NoError(t, err)

	usernames := []string{"zebra", "alpha", "beta"}
	for _, username := range usernames {
		email, _ := valueobjects.NewEmail(username + "@example.com")
		role, _ := valueobjects.NewRole("student")
		passwordHash, _ := utils.HashPassword("Password123!")
		user := entities.NewUser(email, username, role, passwordHash)
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	sortColumn := "username"
	sortDirection := repositories.SortAsc
	limit := 10
	input := usecases.FindUserInput{
		SortColumn:    &sortColumn,
		SortDirection: &sortDirection,
		Limit:         &limit,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.Total)
	assert.Len(t, result.Users, 3)
	assert.Equal(t, "alpha", result.Users[0].Username)
}

func TestFindUser_Integration_CombinedFilters(t *testing.T) {
	db, cleanup, err := integration.SetUpTestDatabase(t, integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	})
	require.NoError(t, err)
	defer cleanup()

	userRepo := integration.SetupUserRepository(db)
	logger := logger.NewNop()

	findUserUC := usecases.NewFindUserUseCase(userRepo, logger)

	ctx := context.Background()

	var _ *sql.DB = db
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%@asto-lms.local'")
	require.NoError(t, err)

	studentRole, _ := valueobjects.NewRole("student")
	activeStatus, _ := valueobjects.NewStatus("active")

	email1, _ := valueobjects.NewEmail("combined1@example.com")
	user1 := entities.NewUser(email1, "combined1", studentRole, "hash")
	user1.Status = activeStatus
	err = userRepo.Create(ctx, user1)
	require.NoError(t, err)

	email2, _ := valueobjects.NewEmail("combined2@example.com")
	user2 := entities.NewUser(email2, "combined2", studentRole, "hash")
	pendingStatus, _ := valueobjects.NewStatus("pending")
	user2.Status = pendingStatus
	err = userRepo.Create(ctx, user2)
	require.NoError(t, err)

	roleFilter := valueobjects.Role("student")
	statusFilter := valueobjects.Status("active")
	limit := 10
	input := usecases.FindUserInput{
		Role:  &roleFilter,
		Status: &statusFilter,
		Limit: &limit,
	}

	result, err := findUserUC.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "combined1@example.com", result.Users[0].Email)
}

