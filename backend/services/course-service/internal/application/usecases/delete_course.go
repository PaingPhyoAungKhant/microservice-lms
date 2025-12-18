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


type DeleteCourseInput struct {
	CourseID string
}

type DeleteCourseOutput struct {
	Message string
}

type DeleteCourseUseCase struct {
	courseRepo repositories.CourseRepository
	publisher  messaging.Publisher
	logger     *logger.Logger
}

func NewDeleteCourseUseCase(
	courseRepo repositories.CourseRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteCourseUseCase {
	return &DeleteCourseUseCase{
		courseRepo: courseRepo,
		publisher:  publisher,
		logger:     logger,
	}
}

func (uc *DeleteCourseUseCase) Execute(ctx context.Context, input DeleteCourseInput) (*DeleteCourseOutput, error) {
	course, err := uc.courseRepo.FindByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	if err := uc.courseRepo.Delete(ctx, input.CourseID); err != nil {
		return nil, fmt.Errorf("failed to delete course: %w", err)
	}

	event := events.CourseDeletedEvent{
		ID:        course.ID,
		DeletedAt: time.Now(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeCourseDeleted, event); err != nil {
			uc.logger.Error("failed to publish course deleted event", zap.Error(err))
		}
	}

	uc.logger.Info("Course deleted successfully",
		zap.String("course_id", course.ID),
	)

	return &DeleteCourseOutput{
		Message: "Course deleted successfully",
	}, nil
}

