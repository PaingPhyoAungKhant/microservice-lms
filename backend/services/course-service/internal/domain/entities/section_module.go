package entities

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentTypeZoom ContentType = "zoom"
)

type ContentStatus string

const (
	ContentStatusDraft   ContentStatus = "draft"
	ContentStatusPending ContentStatus = "pending"
	ContentStatusCreated ContentStatus = "created"
)

type SectionModule struct {
	ID            string
	CourseSectionID string
	ContentID     *string
	Name          string
	Description   string
	ContentType   ContentType
	ContentStatus ContentStatus
	Order         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSectionModule(courseSectionID, name, description string, contentType ContentType, order int) *SectionModule {
	now := time.Now().UTC()
	return &SectionModule{
		ID:              uuid.NewString(),
		CourseSectionID: courseSectionID,
		ContentID:       nil,
		Name:            name,
		Description:     description,
		ContentType:     contentType,
		ContentStatus:   ContentStatusDraft,
		Order:           order,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (sm *SectionModule) Update(name, description string, order int) {
	sm.Name = name
	sm.Description = description
	sm.Order = order
	sm.UpdatedAt = time.Now().UTC()
}

func (sm *SectionModule) UpdateContent(contentID *string, status ContentStatus) {
	sm.ContentID = contentID
	sm.ContentStatus = status
	sm.UpdatedAt = time.Now().UTC()
}

func (sm *SectionModule) UpdateStatus(status ContentStatus) {
	sm.ContentStatus = status
	sm.UpdatedAt = time.Now().UTC()
}

func (sm *SectionModule) UpdateOrder(order int) {
	sm.Order = order
	sm.UpdatedAt = time.Now().UTC()
}

