package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/repositories"
)

type PostgresFileRepository struct {
	db *sql.DB
}

func NewPostgresFileRepository(db *sql.DB) repositories.FileRepository {
	return &PostgresFileRepository{db: db}
}

func (r *PostgresFileRepository) Create(ctx context.Context, file *entities.File) error {
	tagsJSON, err := file.TagsToJSONB()
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO files (
			id, original_filename, stored_filename, bucket_name, mime_type, 
			size_bytes, uploaded_by, tags, created_at, updated_at
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = r.db.ExecContext(ctx, query,
		file.ID,
		file.OriginalFilename,
		file.StoredFilename,
		file.BucketName,
		file.MimeType,
		file.SizeBytes,
		file.UploadedBy,
		tagsJSON,
		file.CreatedAt,
		file.UpdatedAt,
	)
	return err
}

func (r *PostgresFileRepository) FindByID(ctx context.Context, id string) (*entities.File, error) {
	query := `
		SELECT
			id, original_filename, stored_filename, bucket_name, mime_type,
			size_bytes, uploaded_by, tags, created_at, updated_at, deleted_at
		FROM files
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanRowToEntity(row, nil)
}

func (r *PostgresFileRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.File, error) {
	query := `
		SELECT
			id, original_filename, stored_filename, bucket_name, mime_type,
			size_bytes, uploaded_by, tags, created_at, updated_at, deleted_at
		FROM files
		WHERE uploaded_by = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*entities.File
	for rows.Next() {
		file, err := r.scanRowToEntity(nil, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return files, nil
}

func (r *PostgresFileRepository) FindByTags(ctx context.Context, tags []string, limit, offset int) ([]*entities.File, error) {
	query := `
		SELECT
			id, original_filename, stored_filename, bucket_name, mime_type,
			size_bytes, uploaded_by, tags, created_at, updated_at, deleted_at
		FROM files
		WHERE tags ?| $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	tagsArray := fmt.Sprintf("{%s}", strings.Join(tags, ","))
	rows, err := r.db.QueryContext(ctx, query, tagsArray, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*entities.File
	for rows.Next() {
		file, err := r.scanRowToEntity(nil, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return files, nil
}

func (r *PostgresFileRepository) Find(ctx context.Context, query repositories.FileQuery) (*repositories.FileQueryResult, error) {
	var queryBuilder strings.Builder
	whereClauses := []string{"deleted_at IS NULL"}
	whereArgs := []interface{}{}
	argIdx := 1

	if query.UploadedBy != nil && *query.UploadedBy != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("uploaded_by = $%d", argIdx))
		whereArgs = append(whereArgs, *query.UploadedBy)
		argIdx++
	}

	if query.BucketName != nil && *query.BucketName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("bucket_name = $%d", argIdx))
		whereArgs = append(whereArgs, *query.BucketName)
		argIdx++
	}

	if query.MimeType != nil && *query.MimeType != "" {
		if strings.Contains(*query.MimeType, "*") {
			mimePattern := strings.Replace(*query.MimeType, "*", "%", -1)
			whereClauses = append(whereClauses, fmt.Sprintf("mime_type LIKE $%d", argIdx))
			whereArgs = append(whereArgs, mimePattern)
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("mime_type = $%d", argIdx))
			whereArgs = append(whereArgs, *query.MimeType)
		}
		argIdx++
	}

	if len(query.Tags) > 0 {
		tagsArray := fmt.Sprintf("{%s}", strings.Join(query.Tags, ","))
		whereClauses = append(whereClauses, fmt.Sprintf("tags ?| $%d", argIdx))
		whereArgs = append(whereArgs, tagsArray)
		argIdx++
	}

	whereClause := " WHERE " + strings.Join(whereClauses, " AND ")
	args := whereArgs

	queryBuilder.WriteString(`
		SELECT
			id, original_filename, stored_filename, bucket_name, mime_type,
			size_bytes, uploaded_by, tags, created_at, updated_at, deleted_at
		FROM files
	`)
	queryBuilder.WriteString(whereClause)

	// Order by
	orderClause := " ORDER BY created_at DESC"
	if query.SortColumn != nil && *query.SortColumn != "" {
		column := strings.ToLower(*query.SortColumn)
		allowedColumns := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"size_bytes": true,
			"original_filename": true,
		}
		if allowedColumns[column] {
			direction := "DESC"
			if query.SortDirection != nil && *query.SortDirection == "ASC" {
				direction = "ASC"
			}
			orderClause = fmt.Sprintf(" ORDER BY %s %s", column, direction)
		}
	}
	queryBuilder.WriteString(orderClause)

	// Limit and offset
	limit := 20
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
	files := make([]*entities.File, 0, limit)
	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		file, err := r.scanRowToEntity(nil, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	// Count total
	var countBuilder strings.Builder
	countBuilder.WriteString("SELECT COUNT(*) FROM files")
	countBuilder.WriteString(whereClause)

	var totalCount int
	err = r.db.QueryRowContext(ctx, countBuilder.String(), whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &repositories.FileQueryResult{
		Files: files,
		Total: totalCount,
	}, nil
}

func (r *PostgresFileRepository) Update(ctx context.Context, file *entities.File) error {
	tagsJSON, err := file.TagsToJSONB()
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE files
		SET 
			original_filename = $1,
			stored_filename = $2,
			bucket_name = $3,
			mime_type = $4,
			size_bytes = $5,
			tags = $6,
			updated_at = $7,
			deleted_at = $8
		WHERE id = $9
	`

	_, err = r.db.ExecContext(ctx, query,
		file.OriginalFilename,
		file.StoredFilename,
		file.BucketName,
		file.MimeType,
		file.SizeBytes,
		tagsJSON,
		file.UpdatedAt,
		file.DeletedAt,
		file.ID,
	)
	return err
}

func (r *PostgresFileRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (r *PostgresFileRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE files
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete file: %w", err)
	}
	return nil
}

func (r *PostgresFileRepository) scanRowToEntity(row *sql.Row, rows *sql.Rows) (*entities.File, error) {
	var (
		id               string
		originalFilename string
		storedFilename   string
		bucketName       string
		mimeType         string
		sizeBytes        int64
		uploadedBy        string
		tagsJSON         json.RawMessage
		createdAt        time.Time
		updatedAt        time.Time
		deletedAt        sql.NullTime
	)

	var err error
	if row != nil {
		err = row.Scan(&id, &originalFilename, &storedFilename, &bucketName, &mimeType,
			&sizeBytes, &uploadedBy, &tagsJSON, &createdAt, &updatedAt, &deletedAt)
	} else if rows != nil {
		err = rows.Scan(&id, &originalFilename, &storedFilename, &bucketName, &mimeType,
			&sizeBytes, &uploadedBy, &tagsJSON, &createdAt, &updatedAt, &deletedAt)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	tags, err := entities.TagsFromJSONB(tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	var deletedAtPtr *time.Time
	if deletedAt.Valid {
		deletedAtPtr = &deletedAt.Time
	}

	return &entities.File{
		ID:               id,
		OriginalFilename: originalFilename,
		StoredFilename:   storedFilename,
		BucketName:       bucketName,
		MimeType:         mimeType,
		SizeBytes:        sizeBytes,
		UploadedBy:       uploadedBy,
		Tags:             tags,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		DeletedAt:        deletedAtPtr,
	}, nil
}

