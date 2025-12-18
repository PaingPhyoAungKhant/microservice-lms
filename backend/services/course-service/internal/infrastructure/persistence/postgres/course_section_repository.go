package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCourseSectionRepository struct {
	db *sql.DB
}

func NewPostgresCourseSectionRepository(db *sql.DB) repositories.CourseSectionRepository {
	return &PostgresCourseSectionRepository{db: db}
}

func (r *PostgresCourseSectionRepository) Create(ctx context.Context, section *entities.CourseSection) error {
	query := `
		INSERT INTO course_section (id, course_offering_id, name, description, "order", status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		section.ID,
		section.CourseOfferingID,
		section.Name,
		section.Description,
		section.Order,
		section.Status,
		section.CreatedAt,
		section.UpdatedAt,
	)
	return err
}

func (r *PostgresCourseSectionRepository) FindByID(ctx context.Context, id string) (*entities.CourseSection, error) {
	query := `
		SELECT id, course_offering_id, name, description, "order", status, created_at, updated_at
		FROM course_section
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCourseSection(row)
}

func (r *PostgresCourseSectionRepository) FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseSection, error) {
	query := `
		SELECT id, course_offering_id, name, description, "order", status, created_at, updated_at
		FROM course_section
		WHERE course_offering_id = $1
		ORDER BY "order" ASC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, offeringID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections := []*entities.CourseSection{}
	for rows.Next() {
		section, err := r.scanCourseSectionRow(rows)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, nil
}

func (r *PostgresCourseSectionRepository) Update(ctx context.Context, section *entities.CourseSection) error {
	query := `
		UPDATE course_section
		SET name = $1, description = $2, "order" = $3, status = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		section.Name,
		section.Description,
		section.Order,
		section.Status,
		section.UpdatedAt,
		section.ID,
	)
	return err
}

func (r *PostgresCourseSectionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM course_section WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresCourseSectionRepository) DeleteByOfferingID(ctx context.Context, offeringID string) error {
	query := `DELETE FROM course_section WHERE course_offering_id = $1`
	_, err := r.db.ExecContext(ctx, query, offeringID)
	return err
}

func (r *PostgresCourseSectionRepository) scanCourseSection(row *sql.Row) (*entities.CourseSection, error) {
	var section entities.CourseSection
	err := row.Scan(
		&section.ID,
		&section.CourseOfferingID,
		&section.Name,
		&section.Description,
		&section.Order,
		&section.Status,
		&section.CreatedAt,
		&section.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &section, nil
}

func (r *PostgresCourseSectionRepository) scanCourseSectionRow(rows *sql.Rows) (*entities.CourseSection, error) {
	var section entities.CourseSection
	err := rows.Scan(
		&section.ID,
		&section.CourseOfferingID,
		&section.Name,
		&section.Description,
		&section.Order,
		&section.Status,
		&section.CreatedAt,
		&section.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &section, nil
}

