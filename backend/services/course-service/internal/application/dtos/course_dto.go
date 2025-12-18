package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseDTO struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ThumbnailID  *string      `json:"thumbnail_id,omitempty"`
	ThumbnailURL string       `json:"thumbnail_url,omitempty"`
	Categories   []CategoryDTO `json:"categories,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (d *CourseDTO) FromEntity(course *entities.Course, apiGatewayURL string) {
	d.ID = course.ID
	d.Name = course.Name
	d.Description = course.Description
	d.ThumbnailID = course.ThumbnailID
	d.CreatedAt = course.CreatedAt
	d.UpdatedAt = course.UpdatedAt

	// Compute thumbnail URL if thumbnail_id exists
	if course.ThumbnailID != nil && *course.ThumbnailID != "" && apiGatewayURL != "" {
		d.ThumbnailURL = apiGatewayURL + "/api/v1/buckets/course-thumbnails/files/" + *course.ThumbnailID + "/download"
	}
}

type CreateCourseInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	ThumbnailID *string  `json:"thumbnail_id"`
	CategoryIDs []string `json:"category_ids,omitempty"`
}

type UpdateCourseInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	ThumbnailID *string  `json:"thumbnail_id"`
	CategoryIDs []string `json:"category_ids,omitempty"`
}

