package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseSectionDTO struct {
	ID               string    `json:"id"`
	CourseOfferingID string    `json:"course_offering_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Order            int       `json:"order"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (d *CourseSectionDTO) FromEntity(section *entities.CourseSection) {
	d.ID = section.ID
	d.CourseOfferingID = section.CourseOfferingID
	d.Name = section.Name
	d.Description = section.Description
	d.Order = section.Order
	d.Status = string(section.Status)
	d.CreatedAt = section.CreatedAt
	d.UpdatedAt = section.UpdatedAt
}

type CreateCourseSectionInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type UpdateCourseSectionInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Order       int     `json:"order"`
	Status      *string `json:"status,omitempty"`
}

