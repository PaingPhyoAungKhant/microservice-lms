package entities

import (
	"time"

	"github.com/google/uuid"
)

type CourseOfferingInstructor struct {
	ID                 string
	CourseOfferingID   string
	InstructorID       string
	InstructorUsername string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewCourseOfferingInstructor(courseOfferingID, instructorID, instructorUsername string) *CourseOfferingInstructor {
	now := time.Now().UTC()
	return &CourseOfferingInstructor{
		ID:                 uuid.NewString(),
		CourseOfferingID:   courseOfferingID,
		InstructorID:       instructorID,
		InstructorUsername: instructorUsername,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (coi *CourseOfferingInstructor) UpdateInstructorUsername(username string) {
	coi.InstructorUsername = username
}

