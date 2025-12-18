package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseSectionRepository interface {
	Create(ctx context.Context, section *entities.CourseSection) error
	FindByID(ctx context.Context, id string) (*entities.CourseSection, error)
	FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseSection, error)
	Update(ctx context.Context, section *entities.CourseSection) error
	Delete(ctx context.Context, id string) error
	DeleteByOfferingID(ctx context.Context, offeringID string) error
}

