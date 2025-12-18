package integration

import (
	"database/sql"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	cfg := integration.TestDatabaseConfig{
		MigrationPath:    "migrations",
		TablesToCleanUp: []string{"enrollments"},
	}

	db, cleanup, err := integration.SetUpTestDatabase(t, cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	return db, cleanup
}

func SetupEnrollmentRepository(db *sql.DB) repositories.EnrollmentRepository {
	return postgres.NewPostgresEnrollmentRepository(db)
}

