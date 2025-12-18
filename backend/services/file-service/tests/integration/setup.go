package integration

import (
	"database/sql"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	cfg := sharedIntegration.TestDatabaseConfig{
		MigrationPath:    "migrations",
		TablesToCleanUp: []string{"files"},
	}

	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	return db, cleanup
}

func SetupFileRepository(db *sql.DB) repositories.FileRepository {
	return postgres.NewPostgresFileRepository(db)
}

func SetupTestLogger() *logger.Logger {
	return sharedIntegration.SetupTestLogger()
}

