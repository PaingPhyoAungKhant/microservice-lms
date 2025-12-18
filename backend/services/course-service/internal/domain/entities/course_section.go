package entities

import (
	"time"

	"github.com/google/uuid"
)

type SectionStatus string

const (
	SectionStatusDraft     SectionStatus = "draft"
	SectionStatusPublished SectionStatus = "published"
	SectionStatusArchived  SectionStatus = "archived"
)

type CourseSection struct {
	ID               string
	CourseOfferingID string
	Name             string
	Description      string
	Order            int
	Status           SectionStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewCourseSection(courseOfferingID, name, description string, order int) *CourseSection {
	now := time.Now().UTC()
	return &CourseSection{
		ID:               uuid.NewString(),
		CourseOfferingID: courseOfferingID,
		Name:             name,
		Description:      description,
		Order:            order,
		Status:          SectionStatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (cs *CourseSection) Update(name, description string, order int) {
	cs.Name = name
	cs.Description = description
	cs.Order = order
	cs.UpdatedAt = time.Now().UTC()
}

func (cs *CourseSection) UpdateStatus(status SectionStatus) {
	cs.Status = status
	cs.UpdatedAt = time.Now().UTC()
}

func (cs *CourseSection) UpdateOrder(order int) {
	cs.Order = order
	cs.UpdatedAt = time.Now().UTC()
}

