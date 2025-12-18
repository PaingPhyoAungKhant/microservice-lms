package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
)

type PostgresEnrollmentRepository struct {
	db *sql.DB
}

func NewPostgresEnrollmentRepository(db *sql.DB) repositories.EnrollmentRepository {
	return &PostgresEnrollmentRepository{db: db}
}

func (r *PostgresEnrollmentRepository) Create(ctx context.Context, enrollment *entities.Enrollment) error {
	query := `
		INSERT INTO enrollments (
			id, student_id, student_username, course_id, course_name, 
			course_offering_id, course_offering_name, status, created_at, updated_at
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		enrollment.ID,
		enrollment.StudentID,
		enrollment.StudentUsername,
		enrollment.CourseID,
		enrollment.CourseName,
		enrollment.CourseOfferingID,
		enrollment.CourseOfferingName,
		enrollment.Status.String(),
		enrollment.CreatedAt,
		enrollment.UpdatedAt,
	)
	return err
}

func (r *PostgresEnrollmentRepository) FindByID(ctx context.Context, id string) (*entities.Enrollment, error) {
	query := `
		SELECT
			id, student_id, student_username, course_id, course_name,
			course_offering_id, course_offering_name, status, created_at, updated_at
		FROM enrollments
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanRowToEntity(row, nil)
}

func (r *PostgresEnrollmentRepository) Update(ctx context.Context, enrollment *entities.Enrollment) error {
	query := `
		UPDATE enrollments
		SET 
			student_id = $1,
			student_username = $2,
			course_id = $3,
			course_name = $4,
			course_offering_id = $5,
			course_offering_name = $6,
			status = $7,
			updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(
		ctx, query,
		enrollment.StudentID,
		enrollment.StudentUsername,
		enrollment.CourseID,
		enrollment.CourseName,
		enrollment.CourseOfferingID,
		enrollment.CourseOfferingName,
		enrollment.Status.String(),
		enrollment.UpdatedAt,
		enrollment.ID,
	)
	return err
}

func (r *PostgresEnrollmentRepository) UpdateStudentUsername(ctx context.Context, studentID, username string) error {
	query := `
		UPDATE enrollments
		SET student_username = $1, updated_at = CURRENT_TIMESTAMP
		WHERE student_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, username, studentID)
	return err
}

func (r *PostgresEnrollmentRepository) UpdateCourseName(ctx context.Context, courseID, courseName string) error {
	query := `
		UPDATE enrollments
		SET course_name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE course_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, courseName, courseID)
	return err
}

func (r *PostgresEnrollmentRepository) UpdateCourseOfferingName(ctx context.Context, courseOfferingID, courseOfferingName string) error {
	query := `
		UPDATE enrollments
		SET course_offering_name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE course_offering_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, courseOfferingName, courseOfferingID)
	return err
}

func (r *PostgresEnrollmentRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM enrollments
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	return nil
}

func (r *PostgresEnrollmentRepository) buildWhereClause(query repositories.EnrollmentQuery) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argIdx := 1

	// Search query
	if query.SearchQuery != nil && *query.SearchQuery != "" {
		clauses = append(clauses, fmt.Sprintf("(student_username ILIKE $%d OR course_name ILIKE $%d OR course_offering_name ILIKE $%d)", argIdx, argIdx+1, argIdx+2))
		searchPattern := "%" + *query.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIdx += 3
	}

	// Student ID filter
	if query.StudentID != nil && *query.StudentID != "" {
		clauses = append(clauses, fmt.Sprintf("student_id = $%d", argIdx))
		args = append(args, *query.StudentID)
		argIdx++
	}

	// Course ID filter
	if query.CourseID != nil && *query.CourseID != "" {
		clauses = append(clauses, fmt.Sprintf("course_id = $%d", argIdx))
		args = append(args, *query.CourseID)
		argIdx++
	}

	// Course Offering ID filter
	if query.CourseOfferingID != nil && *query.CourseOfferingID != "" {
		clauses = append(clauses, fmt.Sprintf("course_offering_id = $%d", argIdx))
		args = append(args, *query.CourseOfferingID)
		argIdx++
	}

	// Status filter
	if query.Status != nil && *query.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, query.Status.String())
		argIdx++
	}

	whereClause := ""
	if len(clauses) > 0 {
		whereClause = " WHERE " + strings.Join(clauses, " AND ")
	}

	return whereClause, args
}

func (r *PostgresEnrollmentRepository) buildOrderClause(query repositories.EnrollmentQuery) string {
	if query.SortColumn == nil || *query.SortColumn == "" {
		return " ORDER BY created_at DESC"
	}

	validColumns := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"status":     true,
	}

	if !validColumns[*query.SortColumn] {
		return " ORDER BY created_at DESC"
	}

	direction := "DESC"
	if query.SortDirection != nil && *query.SortDirection == repositories.SortAsc {
		direction = "ASC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", *query.SortColumn, direction)
}

func (r *PostgresEnrollmentRepository) Find(ctx context.Context, query repositories.EnrollmentQuery) (*repositories.EnrollmentQueryResult, error) {
	var queryBuilder strings.Builder
	whereClause, whereArgs := r.buildWhereClause(query)
	args := whereArgs
	argIdx := len(args) + 1

	queryBuilder.WriteString(`
		SELECT
			id, student_id, student_username, course_id, course_name,
			course_offering_id, course_offering_name, status, created_at, updated_at
		FROM enrollments
	`)
	queryBuilder.WriteString(whereClause)

	if orderClause := r.buildOrderClause(query); orderClause != "" {
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
	enrollments := make([]*entities.Enrollment, 0, limit)
	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		enrollment, err := r.scanRowToEntity(nil, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	var countBuilder strings.Builder
	countBuilder.WriteString(`SELECT COUNT(*) FROM enrollments`)
	countBuilder.WriteString(whereClause)

	var totalCount int
	err = r.db.QueryRowContext(ctx, countBuilder.String(), whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &repositories.EnrollmentQueryResult{
		Enrollments: enrollments,
		Total:      totalCount,
	}, nil
}

func (r *PostgresEnrollmentRepository) scanRowToEntity(row *sql.Row, rows *sql.Rows) (*entities.Enrollment, error) {
	var (
		id                 string
		studentID          string
		studentUsername    string
		courseID           string
		courseName         string
		courseOfferingID   string
		courseOfferingName string
		status             string
		createdAt          time.Time
		updatedAt          time.Time
	)
	var err error
	if row != nil {
		err = row.Scan(&id, &studentID, &studentUsername, &courseID, &courseName,
			&courseOfferingID, &courseOfferingName, &status, &createdAt, &updatedAt)
	} else if rows != nil {
		err = rows.Scan(&id, &studentID, &studentUsername, &courseID, &courseName,
			&courseOfferingID, &courseOfferingName, &status, &createdAt, &updatedAt)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("enrollment not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	statusVO, err := valueobjects.NewEnrollmentStatus(status)
	if err != nil {
		return nil, fmt.Errorf("invalid enrollment status: %w", err)
	}

	return &entities.Enrollment{
		ID:                 id,
		StudentID:         studentID,
		StudentUsername:   studentUsername,
		CourseID:          courseID,
		CourseName:        courseName,
		CourseOfferingID:  courseOfferingID,
		CourseOfferingName: courseOfferingName,
		Status:            statusVO,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}, nil
}

