package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type CategoryQuery struct {
	SearchQuery   *string
	Limit         *int
	Offset        *int
	SortColumn    *string
	SortDirection *SortDirection
}

type CategoryQueryResult struct {
	Categories []*entities.Category
	Total      int
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	FindByID(ctx context.Context, id string) (*entities.Category, error)
	FindByName(ctx context.Context, name string) (*entities.Category, error)
	Find(ctx context.Context, query CategoryQuery) (*CategoryQueryResult, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id string) error
}

