package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type DeleteCourseSectionInput struct {
	SectionID string
}

type DeleteCourseSectionOutput struct {
	Message string
}

type DeleteCourseSectionUseCase struct {
	sectionRepo repositories.CourseSectionRepository
	publisher   messaging.Publisher
	logger      *logger.Logger
}

func NewDeleteCourseSectionUseCase(
	sectionRepo repositories.CourseSectionRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteCourseSectionUseCase {
	return &DeleteCourseSectionUseCase{
		sectionRepo: sectionRepo,
		publisher:   publisher,
		logger:      logger,
	}
}

func (uc *DeleteCourseSectionUseCase) Execute(ctx context.Context, input DeleteCourseSectionInput) (*DeleteCourseSectionOutput, error) {
	section, err := uc.sectionRepo.FindByID(ctx, input.SectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrCourseSectionNotFound
	}

	if err := uc.sectionRepo.Delete(ctx, input.SectionID); err != nil {
		return nil, fmt.Errorf("failed to delete course section: %w", err)
	}

	event := events.CourseSectionDeletedEvent{
		ID:        section.ID,
		DeletedAt: time.Now(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseSectionDeleted, event); err != nil {
			uc.logger.Error("failed to publish course section deleted event", zap.Error(err))
		}
	}

	uc.logger.Info("Course section deleted successfully",
		zap.String("section_id", section.ID),
	)

	return &DeleteCourseSectionOutput{
		Message: "Course section deleted successfully",
	}, nil
}

