package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) repositories.CategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	query := `
		INSERT INTO category (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)
	return err
}

func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id string) (*entities.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM category
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCategory(row)
}

func (r *PostgresCategoryRepository) FindByName(ctx context.Context, name string) (*entities.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM category
		WHERE name = $1
	`
	row := r.db.QueryRowContext(ctx, query, name)
	return r.scanCategory(row)
}

func (r *PostgresCategoryRepository) Find(ctx context.Context, query repositories.CategoryQuery) (*repositories.CategoryQueryResult, error) {
	whereClause, args := r.buildWhereClause(query)
	argIdx := len(args) + 1

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM category %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	// Select query
	orderBy := "ORDER BY created_at DESC"
	if query.SortColumn != nil && *query.SortColumn != "" {
		direction := "DESC"
		if query.SortDirection != nil && *query.SortDirection == repositories.SortDirectionASC {
			direction = "ASC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", *query.SortColumn, direction)
	}

	limit := 50
	if query.Limit != nil && *query.Limit > 0 {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil && *query.Offset > 0 {
		offset = *query.Offset
	}

	selectQuery := fmt.Sprintf(`
		SELECT id, name, description, created_at, updated_at
		FROM category
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIdx, argIdx+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	categories := []*entities.Category{}
	for rows.Next() {
		category, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return &repositories.CategoryQueryResult{
		Categories: categories,
		Total:      total,
	}, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	query := `
		UPDATE category
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)
	return err
}

func (r *PostgresCategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM category WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCategoryRepository) buildWhereClause(query repositories.CategoryQuery) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if query.SearchQuery != nil && *query.SearchQuery != "" {
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx+1))
		searchPattern := "%" + *query.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
		argIdx += 2
	}

	whereClause := ""
	if len(clauses) > 0 {
		whereClause = "WHERE " + strings.Join(clauses, " AND ")
	}
	return whereClause, args
}

func (r *PostgresCategoryRepository) scanCategory(row *sql.Row) (*entities.Category, error) {
	var category entities.Category
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *PostgresCategoryRepository) scanCategoryRow(rows *sql.Rows) (*entities.Category, error) {
	var category entities.Category
	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

