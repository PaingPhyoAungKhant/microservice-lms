package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CategoryDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *CategoryDTO) FromEntity(category *entities.Category) {
	d.ID = category.ID
	d.Name = category.Name
	d.Description = category.Description
	d.CreatedAt = category.CreatedAt
	d.UpdatedAt = category.UpdatedAt
}

type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

