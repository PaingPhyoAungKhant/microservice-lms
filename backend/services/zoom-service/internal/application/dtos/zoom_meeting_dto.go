package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/domain/entities"
)

type ZoomMeetingDTO struct {
	ID              string     `json:"id"`
	SectionModuleID string     `json:"section_module_id"`
	ZoomMeetingID   string     `json:"zoom_meeting_id"`
	Topic           string     `json:"topic"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	Duration        *int       `json:"duration,omitempty"`
	JoinURL         string     `json:"join_url"`
	StartURL        string     `json:"start_url"`
	Password        *string    `json:"password,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (d *ZoomMeetingDTO) FromEntity(meeting *entities.ZoomMeeting) {
	d.ID = meeting.ID
	d.SectionModuleID = meeting.SectionModuleID
	d.ZoomMeetingID = meeting.ZoomMeetingID
	d.Topic = meeting.Topic
	d.StartTime = meeting.StartTime
	d.Duration = meeting.Duration
	d.JoinURL = meeting.JoinURL
	d.StartURL = meeting.StartURL
	d.Password = meeting.Password
	d.CreatedAt = meeting.CreatedAt
	d.UpdatedAt = meeting.UpdatedAt
}

type CreateZoomMeetingInput struct {
	SectionModuleID string     `json:"section_module_id" binding:"required"`
	Topic           string     `json:"topic" binding:"required"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	Duration        *int       `json:"duration,omitempty"`
	Password        *string    `json:"password,omitempty"`
}

type UpdateZoomMeetingInput struct {
	Topic     string     `json:"topic" binding:"required"`
	StartTime *time.Time `json:"start_time,omitempty"`
	Duration  *int       `json:"duration,omitempty"`
	Password  *string    `json:"password,omitempty"`
}

