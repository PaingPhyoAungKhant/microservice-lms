package events

import "time"

type ZoomMeetingCreatedEvent struct {
	ZoomMeetingID   string    `json:"zoom_meeting_id"`
	SectionModuleID string    `json:"section_module_id"`
	Topic           string    `json:"topic"`
	JoinURL         string    `json:"join_url"`
	StartURL        string    `json:"start_url"`
	CreatedAt       time.Time `json:"created_at"`
}

