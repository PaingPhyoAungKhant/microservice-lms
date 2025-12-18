package repositories

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
)

type SectionModuleRepository interface {
	Create(ctx context.Context, module *entities.SectionModule) error
	FindByID(ctx context.Context, id string) (*entities.SectionModule, error)
	FindBySectionID(ctx context.Context, sectionID string) ([]*entities.SectionModule, error)
	Update(ctx context.Context, module *entities.SectionModule) error
	Delete(ctx context.Context, id string) error
	DeleteBySectionID(ctx context.Context, sectionID string) error
}

