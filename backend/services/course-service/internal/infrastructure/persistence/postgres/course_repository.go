package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCourseRepository struct {
	db *sql.DB
}

func NewPostgresCourseRepository(db *sql.DB) repositories.CourseRepository {
	return &PostgresCourseRepository{db: db}
}

func (r *PostgresCourseRepository) Create(ctx context.Context, course *entities.Course) error {
	query := `
		INSERT INTO course (id, name, description, thumbnail_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.Name,
		course.Description,
		course.ThumbnailID,
		course.CreatedAt,
		course.UpdatedAt,
	)
	return err
}

func (r *PostgresCourseRepository) FindByID(ctx context.Context, id string) (*entities.Course, error) {
	query := `
		SELECT id, name, description, thumbnail_id, created_at, updated_at
		FROM course
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCourse(row)
}

func (r *PostgresCourseRepository) Find(ctx context.Context, query repositories.CourseQuery) (*repositories.CourseQueryResult, error) {
	whereClause, args := r.buildWhereClause(query)
	argIdx := len(args) + 1

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM course %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count courses: %w", err)
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
		SELECT id, name, description, thumbnail_id, created_at, updated_at
		FROM course
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIdx, argIdx+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()

	courses := []*entities.Course{}
	for rows.Next() {
		course, err := r.scanCourseRow(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return &repositories.CourseQueryResult{
		Courses: courses,
		Total:   total,
	}, nil
}

func (r *PostgresCourseRepository) Update(ctx context.Context, course *entities.Course) error {
	query := `
		UPDATE course
		SET name = $1, description = $2, thumbnail_id = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		course.Name,
		course.Description,
		course.ThumbnailID,
		course.UpdatedAt,
		course.ID,
	)
	return err
}

func (r *PostgresCourseRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM course WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCourseRepository) buildWhereClause(query repositories.CourseQuery) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if query.SearchQuery != nil && *query.SearchQuery != "" {
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx+1))
		searchPattern := "%" + *query.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
		argIdx += 2
	}

	if query.CategoryID != nil && *query.CategoryID != "" {
		clauses = append(clauses, fmt.Sprintf(`
			id IN (
				SELECT course_id FROM course_category WHERE category_id = $%d
			)
		`, argIdx))
		args = append(args, *query.CategoryID)
		argIdx++
	}

	whereClause := ""
	if len(clauses) > 0 {
		whereClause = "WHERE " + strings.Join(clauses, " AND ")
	}
	return whereClause, args
}

func (r *PostgresCourseRepository) scanCourse(row *sql.Row) (*entities.Course, error) {
	var course entities.Course
	var thumbnailID sql.NullString
	err := row.Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&thumbnailID,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if thumbnailID.Valid {
		course.ThumbnailID = &thumbnailID.String
	}
	return &course, nil
}

func (r *PostgresCourseRepository) scanCourseRow(rows *sql.Rows) (*entities.Course, error) {
	var course entities.Course
	var thumbnailID sql.NullString
	err := rows.Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&thumbnailID,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if thumbnailID.Valid {
		course.ThumbnailID = &thumbnailID.String
	}
	return &course, nil
}

