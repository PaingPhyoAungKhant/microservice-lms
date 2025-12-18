package entities

import (
	"time"

	"github.com/google/uuid"
)

type ZoomMeeting struct {
	ID              string
	SectionModuleID string
	ZoomMeetingID   string
	Topic           string
	StartTime       *time.Time
	Duration        *int
	JoinURL         string
	StartURL        string
	Password        *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewZoomMeeting(sectionModuleID, zoomMeetingID, topic, joinURL, startURL string, startTime *time.Time, duration *int, password *string) *ZoomMeeting {
	now := time.Now().UTC()
	return &ZoomMeeting{
		ID:              uuid.NewString(),
		SectionModuleID: sectionModuleID,
		ZoomMeetingID:   zoomMeetingID,
		Topic:           topic,
		StartTime:       startTime,
		Duration:        duration,
		JoinURL:         joinURL,
		StartURL:        startURL,
		Password:        password,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (zm *ZoomMeeting) Update(topic string, startTime *time.Time, duration *int, password *string) {
	zm.Topic = topic
	if startTime != nil {
		zm.StartTime = startTime
	}
	if duration != nil {
		zm.Duration = duration
	}
	if password != nil {
		zm.Password = password
	}
	zm.UpdatedAt = time.Now().UTC()
}

func (zm *ZoomMeeting) UpdateURLs(joinURL, startURL string) {
	zm.JoinURL = joinURL
	zm.StartURL = startURL
	zm.UpdatedAt = time.Now().UTC()
}

