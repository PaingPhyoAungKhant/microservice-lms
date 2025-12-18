package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type DeleteEnrollmentInput struct {
	EnrollmentID string
}

type DeleteEnrollmentOutput struct {
	Message string `json:"message"`
}

type DeleteEnrollmentUseCase struct {
	enrollmentRepo repositories.EnrollmentRepository
	publisher       messaging.Publisher
	logger          *logger.Logger
}

func NewDeleteEnrollmentUseCase(
	enrollmentRepo repositories.EnrollmentRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *DeleteEnrollmentUseCase {
	return &DeleteEnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		publisher:      publisher,
		logger:         logger,
	}
}

func (uc *DeleteEnrollmentUseCase) Execute(ctx context.Context, input DeleteEnrollmentInput) (*DeleteEnrollmentOutput, error) {
	enrollment, err := uc.enrollmentRepo.FindByID(ctx, input.EnrollmentID)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}

	if err := uc.enrollmentRepo.Delete(ctx, input.EnrollmentID); err != nil {
		return nil, fmt.Errorf("failed to delete enrollment: %w", err)
	}

	event := events.EnrollmentDeletedEvent{
		ID:                 enrollment.ID,
		StudentID:          enrollment.StudentID,
		StudentUsername:    enrollment.StudentUsername,
		CourseID:           enrollment.CourseID,
		CourseName:         enrollment.CourseName,
		CourseOfferingID:   enrollment.CourseOfferingID,
		CourseOfferingName: enrollment.CourseOfferingName,
		DeletedAt:          time.Now().UTC(),
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeEnrollmentDeleted, event); err != nil {
			uc.logger.Error("Failed to publish enrollment deleted event", zap.Error(err))
		}
	}

	uc.logger.Info("Enrollment Deleted Successfully.",
		zap.String("enrollment_id", enrollment.ID),
		zap.String("student_id", enrollment.StudentID),
		zap.String("course_offering_id", enrollment.CourseOfferingID),
	)

	return &DeleteEnrollmentOutput{
		Message: "Enrollment deleted successfully",
	}, nil
}

