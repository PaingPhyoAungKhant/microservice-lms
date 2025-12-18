package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseOfferingQuery struct {
	SearchQuery   *string
	CourseID      *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *SortDirection
}

type CourseOfferingQueryResult struct {
	Offerings []*entities.CourseOffering
	Total     int
}

type CourseOfferingRepository interface {
	Create(ctx context.Context, offering *entities.CourseOffering) error
	FindByID(ctx context.Context, id string) (*entities.CourseOffering, error)
	FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseOffering, error)
	Find(ctx context.Context, query CourseOfferingQuery) (*CourseOfferingQueryResult, error)
	Update(ctx context.Context, offering *entities.CourseOffering) error
	Delete(ctx context.Context, id string) error
}

