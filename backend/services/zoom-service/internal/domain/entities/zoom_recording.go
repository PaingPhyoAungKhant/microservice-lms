package entities

import (
	"time"

	"github.com/google/uuid"
)

type ZoomRecording struct {
	ID                 string
	ZoomMeetingID     string
	FileID             string
	RecordingType      *string
	RecordingStartTime *time.Time
	RecordingEndTime   *time.Time
	FileSize           *int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewZoomRecording(zoomMeetingID, fileID string, recordingType *string, recordingStartTime, recordingEndTime *time.Time, fileSize *int64) *ZoomRecording {
	now := time.Now().UTC()
	return &ZoomRecording{
		ID:                 uuid.NewString(),
		ZoomMeetingID:     zoomMeetingID,
		FileID:             fileID,
		RecordingType:      recordingType,
		RecordingStartTime: recordingStartTime,
		RecordingEndTime:   recordingEndTime,
		FileSize:           fileSize,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (zr *ZoomRecording) Update(recordingType *string, recordingStartTime, recordingEndTime *time.Time, fileSize *int64) {
	if recordingType != nil {
		zr.RecordingType = recordingType
	}
	if recordingStartTime != nil {
		zr.RecordingStartTime = recordingStartTime
	}
	if recordingEndTime != nil {
		zr.RecordingEndTime = recordingEndTime
	}
	if fileSize != nil {
		zr.FileSize = fileSize
	}
	zr.UpdatedAt = time.Now().UTC()
}

