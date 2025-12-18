package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
)

type FileQuery struct {
	UploadedBy *string
	Tags       []string
	MimeType   *string
	BucketName *string
	Limit      *int
	Offset     *int
	SortColumn *string
	SortDirection *string
}

type FileQueryResult struct {
	Files []*entities.File
	Total int
}

type FileRepository interface {
	Create(ctx context.Context, file *entities.File) error
	FindByID(ctx context.Context, id string) (*entities.File, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.File, error)
	FindByTags(ctx context.Context, tags []string, limit, offset int) ([]*entities.File, error)
	Find(ctx context.Context, query FileQuery) (*FileQueryResult, error)
	Update(ctx context.Context, file *entities.File) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

