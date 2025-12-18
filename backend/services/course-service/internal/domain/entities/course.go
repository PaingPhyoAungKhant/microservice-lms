package entities

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          string
	Name        string
	Description string
	ThumbnailID *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCourse(name, description string, thumbnailID *string) *Course {
	now := time.Now().UTC()
	return &Course{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		ThumbnailID: thumbnailID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *Course) Update(name, description string, thumbnailID *string) {
	c.Name = name
	c.Description = description
	if thumbnailID != nil {
		c.ThumbnailID = thumbnailID
	}
	c.UpdatedAt = time.Now().UTC()
}

func (c *Course) UpdateThumbnail(thumbnailID *string) {
	c.ThumbnailID = thumbnailID
	c.UpdatedAt = time.Now().UTC()
}

