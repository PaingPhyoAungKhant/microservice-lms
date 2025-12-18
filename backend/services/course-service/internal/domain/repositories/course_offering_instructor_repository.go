package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseOfferingInstructorRepository interface {
	Create(ctx context.Context, instructor *entities.CourseOfferingInstructor) error
	FindByOfferingID(ctx context.Context, offeringID string) ([]*entities.CourseOfferingInstructor, error)
	FindByInstructorID(ctx context.Context, instructorID string) ([]*entities.CourseOfferingInstructor, error)
	Delete(ctx context.Context, id string) error
	DeleteByOfferingID(ctx context.Context, offeringID string) error
	Update(ctx context.Context, instructor *entities.CourseOfferingInstructor) error
}

