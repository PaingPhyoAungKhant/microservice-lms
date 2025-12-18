package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
)

type ZoomRecordingDTO struct {
	ID                 string     `json:"id"`
	ZoomMeetingID     string     `json:"zoom_meeting_id"`
	FileID             string     `json:"file_id"`
	RecordingType      *string    `json:"recording_type,omitempty"`
	RecordingStartTime *time.Time `json:"recording_start_time,omitempty"`
	RecordingEndTime   *time.Time `json:"recording_end_time,omitempty"`
	FileSize           *int64     `json:"file_size,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (d *ZoomRecordingDTO) FromEntity(recording *entities.ZoomRecording) {
	d.ID = recording.ID
	d.ZoomMeetingID = recording.ZoomMeetingID
	d.FileID = recording.FileID
	d.RecordingType = recording.RecordingType
	d.RecordingStartTime = recording.RecordingStartTime
	d.RecordingEndTime = recording.RecordingEndTime
	d.FileSize = recording.FileSize
	d.CreatedAt = recording.CreatedAt
	d.UpdatedAt = recording.UpdatedAt
}

type CreateZoomRecordingInput struct {
	ZoomMeetingID     string     `json:"zoom_meeting_id" binding:"required"`
	FileID            string     `json:"file_id" binding:"required"`
	RecordingType     *string    `json:"recording_type,omitempty"`
	RecordingStartTime *time.Time `json:"recording_start_time,omitempty"`
	RecordingEndTime   *time.Time `json:"recording_end_time,omitempty"`
	FileSize          *int64     `json:"file_size,omitempty"`
}

type UpdateZoomRecordingInput struct {
	RecordingType     *string    `json:"recording_type,omitempty"`
	RecordingStartTime *time.Time `json:"recording_start_time,omitempty"`
	RecordingEndTime   *time.Time `json:"recording_end_time,omitempty"`
	FileSize          *int64     `json:"file_size,omitempty"`
}

