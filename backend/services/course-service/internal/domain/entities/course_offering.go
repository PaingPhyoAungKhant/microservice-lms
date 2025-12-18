package entities

import (
	"time"

	"github.com/google/uuid"
)

type OfferingType string

const (
	OfferingTypeOnline   OfferingType = "online"
	OfferingTypeOnCampus OfferingType = "oncampus"
)

type OfferingStatus string

const (
	OfferingStatusPending   OfferingStatus = "pending"
	OfferingStatusActive    OfferingStatus = "active"
	OfferingStatusOngoing   OfferingStatus = "ongoing"
	OfferingStatusCompleted OfferingStatus = "completed"
)

type CourseOffering struct {
	ID             string
	CourseID       string
	Name           string
	Description    string
	OfferingType   OfferingType
	Status         OfferingStatus
	Duration       *string
	ClassTime      *string
	EnrollmentCost float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewCourseOffering(courseID, name, description string, offeringType OfferingType, duration, classTime *string, enrollmentCost float64) *CourseOffering {
	now := time.Now().UTC()
	return &CourseOffering{
		ID:             uuid.NewString(),
		CourseID:       courseID,
		Name:           name,
		Description:    description,
		OfferingType:   offeringType,
		Status:         OfferingStatusPending,
		Duration:       duration,
		ClassTime:      classTime,
		EnrollmentCost: enrollmentCost,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (co *CourseOffering) Update(name, description string, offeringType OfferingType, duration, classTime *string, enrollmentCost float64) {
	co.Name = name
	co.Description = description
	co.OfferingType = offeringType
	if duration != nil {
		co.Duration = duration
	}
	if classTime != nil {
		co.ClassTime = classTime
	}
	co.EnrollmentCost = enrollmentCost
	co.UpdatedAt = time.Now().UTC()
}

func (co *CourseOffering) UpdateStatus(status OfferingStatus) {
	co.Status = status
	co.UpdatedAt = time.Now().UTC()
}

