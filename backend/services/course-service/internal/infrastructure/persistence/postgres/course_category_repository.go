package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresCourseCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCourseCategoryRepository(db *sql.DB) repositories.CourseCategoryRepository {
	return &PostgresCourseCategoryRepository{db: db}
}

func (r *PostgresCourseCategoryRepository) Create(ctx context.Context, courseCategory *entities.CourseCategory) error {
	query := `
		INSERT INTO course_category (id, course_id, category_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (course_id, category_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query,
		courseCategory.ID,
		courseCategory.CourseID,
		courseCategory.CategoryID,
	)
	return err
}

func (r *PostgresCourseCategoryRepository) FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseCategory, error) {
	query := `
		SELECT id, course_id, category_id
		FROM course_category
		WHERE course_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courseCategories := []*entities.CourseCategory{}
	for rows.Next() {
		var courseCategory entities.CourseCategory
		if err := rows.Scan(
			&courseCategory.ID,
			&courseCategory.CourseID,
			&courseCategory.CategoryID,
		); err != nil {
			return nil, err
		}
		courseCategories = append(courseCategories, &courseCategory)
	}

	return courseCategories, nil
}

func (r *PostgresCourseCategoryRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*entities.CourseCategory, error) {
	query := `
		SELECT id, course_id, category_id
		FROM course_category
		WHERE category_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courseCategories := []*entities.CourseCategory{}
	for rows.Next() {
		var courseCategory entities.CourseCategory
		if err := rows.Scan(
			&courseCategory.ID,
			&courseCategory.CourseID,
			&courseCategory.CategoryID,
		); err != nil {
			return nil, err
		}
		courseCategories = append(courseCategories, &courseCategory)
	}

	return courseCategories, nil
}

func (r *PostgresCourseCategoryRepository) Delete(ctx context.Context, courseID, categoryID string) error {
	query := `
		DELETE FROM course_category
		WHERE course_id = $1 AND category_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, courseID, categoryID)
	return err
}

func (r *PostgresCourseCategoryRepository) DeleteByCourseID(ctx context.Context, courseID string) error {
	query := `DELETE FROM course_category WHERE course_id = $1`
	_, err := r.db.ExecContext(ctx, query, courseID)
	return err
}

