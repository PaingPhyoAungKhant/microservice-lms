package integration

import (
	"database/sql"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	sharedIntegration "github.com/paingphyoaungkhant/asto-microservice/shared/testing/integration"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	return SetupTestDB(t)
}

func setupUserRepo(t *testing.T, db *sql.DB) *postgres.PostgresUserRepository {
	repo := postgres.NewPostgresUserRepository(db)
	return repo.(*postgres.PostgresUserRepository)
}

func setupTestLogger() *logger.Logger {
	return sharedIntegration.SetupTestLogger()
}

func setupJwtManager() *utils.JwtManager {
	return sharedIntegration.SetupJwtManager()
}

