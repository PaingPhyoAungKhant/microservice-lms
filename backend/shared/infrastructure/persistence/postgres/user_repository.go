// Package postgres
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email.String(),
		user.Username,
		user.PasswordHash,
		user.Role.String(),
		user.Status.String(),
		user.EmailVerified,
		user.EmailVerifiedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT
			id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanRowToEntity(row, nil)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT
			id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRowContext(ctx, query, email)
	return scanRowToEntity(row, nil)
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
		SELECT id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, username)
	return scanRowToEntity(row, nil)
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET 
			email = $1,
			username = $2,
			password_hash = $3,
			role = $4,
			status = $5,
			email_verified = $6,
			email_verified_at = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(
		ctx, query,
		user.Email.String(),
		user.Username,
		user.PasswordHash,
		user.Role.String(),
		user.Status.String(),
		user.EmailVerified,
		user.EmailVerifiedAt,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func buildWhereClause(query repositories.UserQuery) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argIdx := 1
	whereClause := ""

	// Search query
	if query.SearchQuery != nil && *query.SearchQuery != "" {
		clauses = append(clauses, fmt.Sprintf("(email ILIKE $%d OR username ILIKE $%d)", argIdx, argIdx+1))
		searchPattern := "%" + *query.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
		argIdx += 2
	}

	// Role filter
	if query.Role != nil && *query.Role != "" {
		clauses = append(clauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, query.Role.String())
		argIdx++
	}

	// Status filter
	if query.Status != nil && *query.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, query.Status.String())
		argIdx++
	}

	if len(clauses) > 0 {
		whereClause = " WHERE " + strings.Join(clauses, " AND ")
	}

	return whereClause, args
}

func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (r *PostgresUserRepository) UpdateEmailVerified(ctx context.Context, userID string, verified bool) error {
	var query string
	if verified {
		query = `
			UPDATE users
			SET email_verified = $1, email_verified_at = NOW(), updated_at = NOW()
			WHERE id = $2
		`
	} else {
		query = `
			UPDATE users
			SET email_verified = $1, email_verified_at = NULL, updated_at = NOW()
			WHERE id = $2
		`
	}
	_, err := r.db.ExecContext(ctx, query, verified, userID)
	return err
}

func buildOrderClause(query repositories.UserQuery) string {
	if query.SortColumn == nil || *query.SortColumn == "" {
		return ""
	}
	column := strings.ToLower(*query.SortColumn)
	allowedColumns := map[string]bool{
		"email":      true,
		"username":   true,
		"role":       true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedColumns[column] {
		return ""
	}

	direction := "ASC"
	if query.SortDirection != nil {
		switch *query.SortDirection {
		case repositories.SortAsc:
			direction = "ASC"
		case repositories.SortDesc:
			direction = "DESC"
		default:
			direction = "ASC"
		}
	}

	return fmt.Sprintf(" ORDER BY %s %s", column, direction)
}

func (r *PostgresUserRepository) Find(ctx context.Context, query repositories.UserQuery) (*repositories.UserQueryResult, error) {
	var queryBuilder strings.Builder
	whereClause, whereArgs := buildWhereClause(query)
	args := whereArgs
	argIdx := len(args) + 1

	queryBuilder.WriteString(`
		SELECT
			id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at
		FROM users
	`)
	queryBuilder.WriteString(whereClause)

	if orderClause := buildOrderClause(query); orderClause != "" {
		queryBuilder.WriteString(orderClause)
	}

	limit := 10
	offset := 0
	if query.Limit != nil && *query.Limit > 0 {
		limit = *query.Limit
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argIdx))
		args = append(args, limit)
		argIdx++
	}

	if query.Offset != nil && *query.Offset > 0 {
		offset = *query.Offset
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argIdx))
		args = append(args, offset)
		argIdx++
	}

	finalQuery := queryBuilder.String()
	users := make([]*entities.User, 0, limit)
	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user, err := scanRowToEntity(nil, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	var countBuilder strings.Builder
	countBuilder.WriteString(`
		SELECT COUNT(*) FROM users
	`)
	countBuilder.WriteString(whereClause)

	var totalCount int
	err = r.db.QueryRowContext(ctx, countBuilder.String(), whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &repositories.UserQueryResult{
		Users: users,
		Total: totalCount,
	}, nil
}

func scanRowToEntity(row *sql.Row, rows *sql.Rows) (*entities.User, error) {
	var (
		userID         string
		email          string
		username       string
		passwordHash   string
		role           string
		status         string
		emailVerified  bool
		emailVerifiedAt sql.NullTime
		createdAt      time.Time
		updatedAt      time.Time
	)
	var err error
	if row != nil {
		err = row.Scan(&userID, &email, &username, &passwordHash, &role, &status, &emailVerified, &emailVerifiedAt, &createdAt, &updatedAt)
	} else if rows != nil {
		err = rows.Scan(&userID, &email, &username, &passwordHash, &role, &status, &emailVerified, &emailVerifiedAt, &createdAt, &updatedAt)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	emailVO, _ := valueobjects.NewEmail(email)
	roleVO, _ := valueobjects.NewRole(role)
	statusVO, _ := valueobjects.NewStatus(status)

	var verifiedAt *time.Time
	if emailVerifiedAt.Valid {
		verifiedAt = &emailVerifiedAt.Time
	}

	return &entities.User{
		ID:             userID,
		Email:          emailVO,
		Username:       username,
		PasswordHash:   passwordHash,
		Role:           roleVO,
		Status:         statusVO,
		EmailVerified:  emailVerified,
		EmailVerifiedAt: verifiedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

