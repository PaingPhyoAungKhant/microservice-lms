package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CourseQuery struct {
	SearchQuery   *string
	CategoryID    *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *SortDirection
}

type CourseQueryResult struct {
	Courses []*entities.Course
	Total   int
}

type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	FindByID(ctx context.Context, id string) (*entities.Course, error)
	Find(ctx context.Context, query CourseQuery) (*CourseQueryResult, error)
	Update(ctx context.Context, course *entities.Course) error
	Delete(ctx context.Context, id string) error
}

type CourseCategoryRepository interface {
	Create(ctx context.Context, courseCategory *entities.CourseCategory) error
	FindByCourseID(ctx context.Context, courseID string) ([]*entities.CourseCategory, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*entities.CourseCategory, error)
	Delete(ctx context.Context, courseID, categoryID string) error
	DeleteByCourseID(ctx context.Context, courseID string) error
}

type SortDirection string

const (
	SortDirectionASC  SortDirection = "ASC"
	SortDirectionDESC SortDirection = "DESC"
)

