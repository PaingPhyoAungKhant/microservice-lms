package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type SectionModuleDTO struct {
	ID              string     `json:"id"`
	CourseSectionID string     `json:"course_section_id"`
	ContentID       *string    `json:"content_id,omitempty"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ContentType     string     `json:"content_type"`
	ContentStatus   string     `json:"content_status"`
	Order           int        `json:"order"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (d *SectionModuleDTO) FromEntity(module *entities.SectionModule) {
	d.ID = module.ID
	d.CourseSectionID = module.CourseSectionID
	d.ContentID = module.ContentID
	d.Name = module.Name
	d.Description = module.Description
	d.ContentType = string(module.ContentType)
	d.ContentStatus = string(module.ContentStatus)
	d.Order = module.Order
	d.CreatedAt = module.CreatedAt
	d.UpdatedAt = module.UpdatedAt
}

type CreateSectionModuleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ContentType string `json:"content_type" binding:"required"`
	Order       int    `json:"order"`
}

type UpdateSectionModuleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

