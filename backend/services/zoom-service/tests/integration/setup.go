package integration

import (
	"database/sql"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	cfg := sharedIntegration.TestDatabaseConfig{
		MigrationPath:    "migrations",
		TablesToCleanUp: []string{"zoom_recording", "zoom_meeting"},
	}

	db, cleanup, err := sharedIntegration.SetUpTestDatabase(t, cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	return db, cleanup
}

func SetupZoomMeetingRepository(db *sql.DB) repositories.ZoomMeetingRepository {
	return postgres.NewPostgresZoomMeetingRepository(db)
}

func SetupZoomRecordingRepository(db *sql.DB) repositories.ZoomRecordingRepository {
	return postgres.NewPostgresZoomRecordingRepository(db)
}

func SetupTestLogger() *logger.Logger {
	return sharedIntegration.SetupTestLogger()
}

