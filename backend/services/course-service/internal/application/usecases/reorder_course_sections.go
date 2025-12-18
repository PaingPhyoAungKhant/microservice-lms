package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

var (
	ErrInvalidReorderInput = errors.New("invalid reorder input")
)

type ReorderItem struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

type ReorderCourseSectionsInput struct {
	OfferingID string        `json:"offering_id"`
	Items      []ReorderItem `json:"items"`
}

type ReorderCourseSectionsUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	logger      *logger.Logger
}

func NewReorderCourseSectionsUseCase(
	sectionRepo repositories.CourseSectionRepository,
	logger *logger.Logger,
) *ReorderCourseSectionsUseCase {
	return &ReorderCourseSectionsUseCase{
		sectionRepo: sectionRepo,
		logger:      logger,
	}
}

func (uc *ReorderCourseSectionsUseCase) Execute(ctx context.Context, input ReorderCourseSectionsInput) error {
	if input.OfferingID == "" {
		return ErrInvalidReorderInput
	}
	if len(input.Items) == 0 {
		return ErrInvalidReorderInput
	}

	for _, item := range input.Items {
		section, err := uc.sectionRepo.FindByID(ctx, item.ID)
		if err != nil {
			return err
		}
		if section == nil {
			return ErrCourseSectionNotFound
		}
		if section.CourseOfferingID != input.OfferingID {
			return ErrInvalidReorderInput
		}
	}

	for _, item := range input.Items {
		section, err := uc.sectionRepo.FindByID(ctx, item.ID)
		if err != nil {
			return err
		}
		if section == nil {
			continue
		}

		section.Order = item.Order
		section.UpdatedAt = time.Now().UTC()

		if err := uc.sectionRepo.Update(ctx, section); err != nil {
			uc.logger.Error("Failed to update section order", zap.String("section_id", item.ID), zap.Error(err))
			return err
		}
	}

	uc.logger.Info("Course sections reordered successfully", zap.String("offering_id", input.OfferingID))
	return nil
}

