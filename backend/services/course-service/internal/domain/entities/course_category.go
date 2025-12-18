package entities

import (
	"github.com/google/uuid"
)

type CourseCategory struct {
	ID         string
	CourseID   string
	CategoryID string
}

func NewCourseCategory(courseID, categoryID string) *CourseCategory {
	return &CourseCategory{
		ID:         uuid.NewString(),
		CourseID:   courseID,
		CategoryID: categoryID,
	}
}

