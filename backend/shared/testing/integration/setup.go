package integration

import (
	"database/sql"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
)

func SetupUserRepository(db *sql.DB) repositories.UserRepository {
	return postgres.NewPostgresUserRepository(db)
}

func SetupTestLogger() *logger.Logger {
	return logger.NewNop()
}

func SetupJwtManager() *utils.JwtManager {
	return utils.NewJwtManager(
		"test-secret-key",
		15*time.Minute,
		24*time.Hour,
	)
}

