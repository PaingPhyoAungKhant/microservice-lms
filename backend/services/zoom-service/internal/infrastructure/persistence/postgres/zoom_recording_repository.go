package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type PostgresZoomRecordingRepository struct {
	db *sql.DB
}

func NewPostgresZoomRecordingRepository(db *sql.DB) repositories.ZoomRecordingRepository {
	return &PostgresZoomRecordingRepository{db: db}
}

func (r *PostgresZoomRecordingRepository) Create(ctx context.Context, recording *entities.ZoomRecording) error {
	query := `
		INSERT INTO zoom_recording (id, zoom_meeting_id, file_id, recording_type, recording_start_time, recording_end_time, file_size, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		recording.ID,
		recording.ZoomMeetingID,
		recording.FileID,
		recording.RecordingType,
		recording.RecordingStartTime,
		recording.RecordingEndTime,
		recording.FileSize,
		recording.CreatedAt,
		recording.UpdatedAt,
	)
	return err
}

func (r *PostgresZoomRecordingRepository) FindByID(ctx context.Context, id string) (*entities.ZoomRecording, error) {
	query := `
		SELECT id, zoom_meeting_id, file_id, recording_type, recording_start_time, recording_end_time, file_size, created_at, updated_at
		FROM zoom_recording
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanZoomRecording(row)
}

func (r *PostgresZoomRecordingRepository) FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) ([]*entities.ZoomRecording, error) {
	query := `
		SELECT id, zoom_meeting_id, file_id, recording_type, recording_start_time, recording_end_time, file_size, created_at, updated_at
		FROM zoom_recording
		WHERE zoom_meeting_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, zoomMeetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recordings := []*entities.ZoomRecording{}
	for rows.Next() {
		recording, err := r.scanZoomRecordingRow(rows)
		if err != nil {
			return nil, err
		}
		recordings = append(recordings, recording)
	}

	return recordings, nil
}

func (r *PostgresZoomRecordingRepository) Update(ctx context.Context, recording *entities.ZoomRecording) error {
	query := `
		UPDATE zoom_recording
		SET recording_type = $1, recording_start_time = $2, recording_end_time = $3, file_size = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		recording.RecordingType,
		recording.RecordingStartTime,
		recording.RecordingEndTime,
		recording.FileSize,
		recording.UpdatedAt,
		recording.ID,
	)
	return err
}

func (r *PostgresZoomRecordingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM zoom_recording WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresZoomRecordingRepository) scanZoomRecording(row *sql.Row) (*entities.ZoomRecording, error) {
	var recording entities.ZoomRecording
	var recordingType sql.NullString
	var recordingStartTime, recordingEndTime sql.NullTime
	var fileSize sql.NullInt64
	err := row.Scan(
		&recording.ID,
		&recording.ZoomMeetingID,
		&recording.FileID,
		&recordingType,
		&recordingStartTime,
		&recordingEndTime,
		&fileSize,
		&recording.CreatedAt,
		&recording.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if recordingType.Valid {
		recording.RecordingType = &recordingType.String
	}
	if recordingStartTime.Valid {
		recording.RecordingStartTime = &recordingStartTime.Time
	}
	if recordingEndTime.Valid {
		recording.RecordingEndTime = &recordingEndTime.Time
	}
	if fileSize.Valid {
		recording.FileSize = &fileSize.Int64
	}
	return &recording, nil
}

func (r *PostgresZoomRecordingRepository) scanZoomRecordingRow(rows *sql.Rows) (*entities.ZoomRecording, error) {
	var recording entities.ZoomRecording
	var recordingType sql.NullString
	var recordingStartTime, recordingEndTime sql.NullTime
	var fileSize sql.NullInt64
	err := rows.Scan(
		&recording.ID,
		&recording.ZoomMeetingID,
		&recording.FileID,
		&recordingType,
		&recordingStartTime,
		&recordingEndTime,
		&fileSize,
		&recording.CreatedAt,
		&recording.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if recordingType.Valid {
		recording.RecordingType = &recordingType.String
	}
	if recordingStartTime.Valid {
		recording.RecordingStartTime = &recordingStartTime.Time
	}
	if recordingEndTime.Valid {
		recording.RecordingEndTime = &recordingEndTime.Time
	}
	if fileSize.Valid {
		recording.FileSize = &fileSize.Int64
	}
	return &recording, nil
}

