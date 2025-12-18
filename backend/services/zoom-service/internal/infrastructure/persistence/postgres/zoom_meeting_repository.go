package postgres

import (
	"context"
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/repositories"
)

type PostgresZoomMeetingRepository struct {
	db *sql.DB
}

func NewPostgresZoomMeetingRepository(db *sql.DB) repositories.ZoomMeetingRepository {
	return &PostgresZoomMeetingRepository{db: db}
}

func (r *PostgresZoomMeetingRepository) Create(ctx context.Context, meeting *entities.ZoomMeeting) error {
	query := `
		INSERT INTO zoom_meeting (id, section_module_id, zoom_meeting_id, topic, start_time, duration, join_url, start_url, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		meeting.ID,
		meeting.SectionModuleID,
		meeting.ZoomMeetingID,
		meeting.Topic,
		meeting.StartTime,
		meeting.Duration,
		meeting.JoinURL,
		meeting.StartURL,
		meeting.Password,
		meeting.CreatedAt,
		meeting.UpdatedAt,
	)
	return err
}

func (r *PostgresZoomMeetingRepository) FindByID(ctx context.Context, id string) (*entities.ZoomMeeting, error) {
	query := `
		SELECT id, section_module_id, zoom_meeting_id, topic, start_time, duration, join_url, start_url, password, created_at, updated_at
		FROM zoom_meeting
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanZoomMeeting(row)
}

func (r *PostgresZoomMeetingRepository) FindByZoomMeetingID(ctx context.Context, zoomMeetingID string) (*entities.ZoomMeeting, error) {
	query := `
		SELECT id, section_module_id, zoom_meeting_id, topic, start_time, duration, join_url, start_url, password, created_at, updated_at
		FROM zoom_meeting
		WHERE zoom_meeting_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, zoomMeetingID)
	return r.scanZoomMeeting(row)
}

func (r *PostgresZoomMeetingRepository) FindBySectionModuleID(ctx context.Context, sectionModuleID string) (*entities.ZoomMeeting, error) {
	query := `
		SELECT id, section_module_id, zoom_meeting_id, topic, start_time, duration, join_url, start_url, password, created_at, updated_at
		FROM zoom_meeting
		WHERE section_module_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, sectionModuleID)
	return r.scanZoomMeeting(row)
}

func (r *PostgresZoomMeetingRepository) Update(ctx context.Context, meeting *entities.ZoomMeeting) error {
	query := `
		UPDATE zoom_meeting
		SET topic = $1, start_time = $2, duration = $3, join_url = $4, start_url = $5, password = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		meeting.Topic,
		meeting.StartTime,
		meeting.Duration,
		meeting.JoinURL,
		meeting.StartURL,
		meeting.Password,
		meeting.UpdatedAt,
		meeting.ID,
	)
	return err
}

func (r *PostgresZoomMeetingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM zoom_meeting WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresZoomMeetingRepository) scanZoomMeeting(row *sql.Row) (*entities.ZoomMeeting, error) {
	var meeting entities.ZoomMeeting
	var startTime sql.NullTime
	var duration sql.NullInt32
	var password sql.NullString
	err := row.Scan(
		&meeting.ID,
		&meeting.SectionModuleID,
		&meeting.ZoomMeetingID,
		&meeting.Topic,
		&startTime,
		&duration,
		&meeting.JoinURL,
		&meeting.StartURL,
		&password,
		&meeting.CreatedAt,
		&meeting.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if startTime.Valid {
		meeting.StartTime = &startTime.Time
	}
	if duration.Valid {
		d := int(duration.Int32)
		meeting.Duration = &d
	}
	if password.Valid {
		meeting.Password = &password.String
	}
	return &meeting, nil
}

