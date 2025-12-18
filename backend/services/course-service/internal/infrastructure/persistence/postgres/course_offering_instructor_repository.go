package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCourseOfferingInstructorRepository struct {
	db *sql.DB
}

func NewPostgresCourseOfferingInstructorRepository(db *sql.DB) repositories.CourseOfferingInstructorRepository {
	return &PostgresCourseOfferingInstructorRepository{db: db}
}

func (r *PostgresCourseOfferingInstructorRepository) Create(ctx context.Context, instructor *entities.CourseOfferingInstructor) error {
	query := `
		INSERT INTO course_offering_instructor (id, course_offering_id, instructor_id, instructor_username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (course_offering_id, instructor_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query,
		instructor.ID,
		instructor.CourseOfferingID,
		instructor.InstructorID,
		instructor.InstructorUsername,
		instructor.CreatedAt,
		instructor.UpdatedAt,
	)
	return err
}

func (r *PostgresCourseOfferingInstructorRepository) FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseOfferingInstructor, error) {
	query := `
		SELECT id, course_offering_id, instructor_id, instructor_username, created_at, updated_at
		FROM course_offering_instructor
		WHERE course_offering_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, offeringID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instructors := []*entities.CourseOfferingInstructor{}
	for rows.Next() {
		instructor, err := r.scanCourseOfferingInstructor(rows)
		if err != nil {
			return nil, err
		}
		instructors = append(instructors, instructor)
	}

	return instructors, nil
}

func (r *PostgresCourseOfferingInstructorRepository) FindByInstructorID(ctx context.Context, instructorID string) ([]*entities.CourseOfferingInstructor, error) {
	query := `
		SELECT id, course_offering_id, instructor_id, instructor_username, created_at, updated_at
		FROM course_offering_instructor
		WHERE instructor_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instructors := []*entities.CourseOfferingInstructor{}
	for rows.Next() {
		instructor, err := r.scanCourseOfferingInstructor(rows)
		if err != nil {
			return nil, err
		}
		instructors = append(instructors, instructor)
	}

	return instructors, nil
}

func (r *PostgresCourseOfferingInstructorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM course_offering_instructor WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCourseOfferingInstructorRepository) DeleteByOfferingID(ctx context.Context, offeringID string) error {
	query := `DELETE FROM course_offering_instructor WHERE course_offering_id = $1`
	_, err := r.db.ExecContext(ctx, query, offeringID)
	return err
}

func (r *PostgresCourseOfferingInstructorRepository) Update(ctx context.Context, instructor *entities.CourseOfferingInstructor) error {
	query := `
		UPDATE course_offering_instructor
		SET instructor_username = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query,
		instructor.InstructorUsername,
		instructor.UpdatedAt,
		instructor.ID,
	)
	return err
}

func (r *PostgresCourseOfferingInstructorRepository) scanCourseOfferingInstructor(rows *sql.Rows) (*entities.CourseOfferingInstructor, error) {
	var instructor entities.CourseOfferingInstructor
	err := rows.Scan(
		&instructor.ID,
		&instructor.CourseOfferingID,
		&instructor.InstructorID,
		&instructor.InstructorUsername,
		&instructor.CreatedAt,
		&instructor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &instructor, nil
}

