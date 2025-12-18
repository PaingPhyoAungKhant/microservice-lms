package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCourseOfferingRepository struct {
	db *sql.DB
}

func NewPostgresCourseOfferingRepository(db *sql.DB) repositories.CourseOfferingRepository {
	return &PostgresCourseOfferingRepository{db: db}
}

func (r *PostgresCourseOfferingRepository) Create(ctx context.Context, offering *entities.CourseOffering) error {
	query := `
		INSERT INTO course_offering (id, course_id, name, description, offering_type, status, duration, class_time, enrollment_cost, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		offering.ID,
		offering.CourseID,
		offering.Name,
		offering.Description,
		offering.OfferingType,
		offering.Status,
		offering.Duration,
		offering.ClassTime,
		offering.EnrollmentCost,
		offering.CreatedAt,
		offering.UpdatedAt,
	)
	return err
}

func (r *PostgresCourseOfferingRepository) FindByID(ctx context.Context, id string) (*entities.CourseOffering, error) {
	query := `
		SELECT id, course_id, name, description, offering_type, status, duration, class_time, enrollment_cost, created_at, updated_at
		FROM course_offering
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCourseOffering(row)
}

func (r *PostgresCourseOfferingRepository) FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseOffering, error) {
	query := `
		SELECT id, course_id, name, description, offering_type, status, duration, class_time, enrollment_cost, created_at, updated_at
		FROM course_offering
		WHERE course_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offerings := []*entities.CourseOffering{}
	for rows.Next() {
		offering, err := r.scanCourseOfferingRow(rows)
		if err != nil {
			return nil, err
		}
		offerings = append(offerings, offering)
	}

	return offerings, nil
}

func (r *PostgresCourseOfferingRepository) Find(ctx context.Context, query repositories.CourseOfferingQuery) (*repositories.CourseOfferingQueryResult, error) {
	whereClause, args := r.buildWhereClause(query)
	argIdx := len(args) + 1

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM course_offering %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count course offerings: %w", err)
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
		SELECT id, course_id, name, description, offering_type, status, duration, class_time, enrollment_cost, created_at, updated_at
		FROM course_offering
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIdx, argIdx+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query course offerings: %w", err)
	}
	defer rows.Close()

	offerings := []*entities.CourseOffering{}
	for rows.Next() {
		offering, err := r.scanCourseOfferingRow(rows)
		if err != nil {
			return nil, err
		}
		offerings = append(offerings, offering)
	}

	return &repositories.CourseOfferingQueryResult{
		Offerings: offerings,
		Total:     total,
	}, nil
}

func (r *PostgresCourseOfferingRepository) buildWhereClause(query repositories.CourseOfferingQuery) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if query.SearchQuery != nil && *query.SearchQuery != "" {
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx+1))
		searchPattern := "%" + *query.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
		argIdx += 2
	}

	if query.CourseID != nil && *query.CourseID != "" {
		clauses = append(clauses, fmt.Sprintf("course_id = $%d", argIdx))
		args = append(args, *query.CourseID)
		argIdx++
	}

	whereClause := ""
	if len(clauses) > 0 {
		whereClause = "WHERE " + strings.Join(clauses, " AND ")
	}
	return whereClause, args
}

func (r *PostgresCourseOfferingRepository) Update(ctx context.Context, offering *entities.CourseOffering) error {
	query := `
		UPDATE course_offering
		SET name = $1, description = $2, offering_type = $3, status = $4, duration = $5, class_time = $6, enrollment_cost = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		offering.Name,
		offering.Description,
		offering.OfferingType,
		offering.Status,
		offering.Duration,
		offering.ClassTime,
		offering.EnrollmentCost,
		offering.UpdatedAt,
		offering.ID,
	)
	return err
}

func (r *PostgresCourseOfferingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM course_offering WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCourseOfferingRepository) scanCourseOffering(row *sql.Row) (*entities.CourseOffering, error) {
	var offering entities.CourseOffering
	var duration, classTime sql.NullString
	err := row.Scan(
		&offering.ID,
		&offering.CourseID,
		&offering.Name,
		&offering.Description,
		&offering.OfferingType,
		&offering.Status,
		&duration,
		&classTime,
		&offering.EnrollmentCost,
		&offering.CreatedAt,
		&offering.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if duration.Valid {
		offering.Duration = &duration.String
	}
	if classTime.Valid {
		offering.ClassTime = &classTime.String
	}
	return &offering, nil
}

func (r *PostgresCourseOfferingRepository) scanCourseOfferingRow(rows *sql.Rows) (*entities.CourseOffering, error) {
	var offering entities.CourseOffering
	var duration, classTime sql.NullString
	err := rows.Scan(
		&offering.ID,
		&offering.CourseID,
		&offering.Name,
		&offering.Description,
		&offering.OfferingType,
		&offering.Status,
		&duration,
		&classTime,
		&offering.EnrollmentCost,
		&offering.CreatedAt,
		&offering.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if duration.Valid {
		offering.Duration = &duration.String
	}
	if classTime.Valid {
		offering.ClassTime = &classTime.String
	}
	return &offering, nil
}

