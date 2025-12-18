package integration

import (
	"database/sql"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	cfg := integration.TestDatabaseConfig{
		MigrationPath:  "migrations",
		TablesToCleanUp: []string{"users"},
	}

	db, cleanup, err := integration.SetUpTestDatabase(t, cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	return db, cleanup
}

