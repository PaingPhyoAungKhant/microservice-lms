package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
)

type PostgresSectionModuleRepository struct {
	db *sql.DB
}

func NewPostgresSectionModuleRepository(db *sql.DB) repositories.SectionModuleRepository {
	return &PostgresSectionModuleRepository{db: db}
}

func (r *PostgresSectionModuleRepository) Create(ctx context.Context, module *entities.SectionModule) error {
	query := `
		INSERT INTO section_module (id, course_section_id, content_id, name, description, content_type, content_status, "order", created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		module.ID,
		module.CourseSectionID,
		module.ContentID,
		module.Name,
		module.Description,
		module.ContentType,
		module.ContentStatus,
		module.Order,
		module.CreatedAt,
		module.UpdatedAt,
	)
	return err
}

func (r *PostgresSectionModuleRepository) FindByID(ctx context.Context, id string) (*entities.SectionModule, error) {
	query := `
		SELECT id, course_section_id, content_id, name, description, content_type, content_status, "order", created_at, updated_at
		FROM section_module
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanSectionModule(row)
}

func (r *PostgresSectionModuleRepository) FindBySectionID(ctx context.Context, sectionID string) ([]*entities.SectionModule, error) {
	query := `
		SELECT id, course_section_id, content_id, name, description, content_type, content_status, "order", created_at, updated_at
		FROM section_module
		WHERE course_section_id = $1
		ORDER BY "order" ASC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := []*entities.SectionModule{}
	for rows.Next() {
		module, err := r.scanSectionModuleRow(rows)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func (r *PostgresSectionModuleRepository) Update(ctx context.Context, module *entities.SectionModule) error {
	query := `
		UPDATE section_module
		SET name = $1, description = $2, content_id = $3, content_type = $4, content_status = $5, "order" = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		module.Name,
		module.Description,
		module.ContentID,
		module.ContentType,
		module.ContentStatus,
		module.Order,
		module.UpdatedAt,
		module.ID,
	)
	return err
}

func (r *PostgresSectionModuleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM section_module WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresSectionModuleRepository) DeleteBySectionID(ctx context.Context, sectionID string) error {
	query := `DELETE FROM section_module WHERE course_section_id = $1`
	_, err := r.db.ExecContext(ctx, query, sectionID)
	return err
}

func (r *PostgresSectionModuleRepository) scanSectionModule(row *sql.Row) (*entities.SectionModule, error) {
	var module entities.SectionModule
	var contentID sql.NullString
	err := row.Scan(
		&module.ID,
		&module.CourseSectionID,
		&contentID,
		&module.Name,
		&module.Description,
		&module.ContentType,
		&module.ContentStatus,
		&module.Order,
		&module.CreatedAt,
		&module.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if contentID.Valid {
		module.ContentID = &contentID.String
	}
	return &module, nil
}

func (r *PostgresSectionModuleRepository) scanSectionModuleRow(rows *sql.Rows) (*entities.SectionModule, error) {
	var module entities.SectionModule
	var contentID sql.NullString
	err := rows.Scan(
		&module.ID,
		&module.CourseSectionID,
		&contentID,
		&module.Name,
		&module.Description,
		&module.ContentType,
		&module.ContentStatus,
		&module.Order,
		&module.CreatedAt,
		&module.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if contentID.Valid {
		module.ContentID = &contentID.String
	}
	return &module, nil
}

