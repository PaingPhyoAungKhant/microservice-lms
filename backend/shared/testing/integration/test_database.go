package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	testPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDatabaseConfig struct {
	DBName string
	Username string
	Password string
	TablesToCleanUp []string
	MigrationPath string
}

func SetUpTestDatabase(t *testing.T, config TestDatabaseConfig) (*sql.DB, func(), error) {
	t.Helper()
  
	if config.DBName == "" {
		config.DBName = "testdb"
	}
	if config.Username == "" {
		config.Username = "postgres"
	}
	if config.Password == "" {
		config.Password = "postgres"
	}
	if config.TablesToCleanUp == nil {
		config.TablesToCleanUp = []string{}
	}
	if config.MigrationPath == "" {
		config.MigrationPath = "migrations"
	}
	ctx := context.Background()

	pgContainer, err := testPostgres.Run(ctx, 
	 "postgres:16-alpine",
	 testPostgres.WithDatabase(config.DBName),
    testPostgres.WithUsername(config.Username),
    testPostgres.WithPassword(config.Password),
    testPostgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}
	connStr, err := pgContainer.ConnectionString(ctx,"sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := runMigrations(db, config.MigrationPath); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	cleanup := func() {
		cleanUpTestDatabase(db, config.TablesToCleanUp)
		db.Close()
		pgContainer.Terminate(ctx)
	}
	return db, cleanup, nil
}

func resolveMigrationPath(migrationPath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	if filepath.IsAbs(migrationPath) {
		if _, err := os.Stat(migrationPath); err == nil {
			return migrationPath, nil
		}
		return "", fmt.Errorf("migrations directory does not exist: %s", migrationPath)
	}

	resolvedPath := filepath.Join(wd, migrationPath)
	if absPath, err := filepath.Abs(resolvedPath); err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}

	currentDir := wd
	for i := 0; i < 10; i++ {
		candidate := filepath.Join(currentDir, migrationPath)
		if absCandidate, err := filepath.Abs(candidate); err == nil {
			if _, err := os.Stat(absCandidate); err == nil {
				return absCandidate, nil
			}
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("migrations directory not found: %s (searched from: %s)", migrationPath, wd)
	}

func runMigrations(db *sql.DB, migrationPath string) error {
	migrationsAbsDir, err := resolveMigrationPath(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to resolve migration path: %w", err)
	}

	migrationsURL := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsAbsDir))

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance (URL: %s): %w", migrationsURL, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate database (dir: %s): %w", migrationsAbsDir, err)
	}

	return nil
}

func cleanUpTestDatabase(db *sql.DB, tablesToCleanUp []string) error {
	for _, table := range tablesToCleanUp {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to clean up table %s: %w", table, err)
		}
	}
	return nil
}