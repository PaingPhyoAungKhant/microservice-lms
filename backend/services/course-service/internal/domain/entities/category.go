package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCategory(name, description string) *Category {
	now := time.Now().UTC()
	return &Category{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *Category) Update(name, description string) {
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now().UTC()
}

