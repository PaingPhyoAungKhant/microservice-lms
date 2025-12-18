package usecases

import (
	"context"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type UpdateEnrollmentStatusInput struct {
	EnrollmentID string
	Status       valueobjects.EnrollmentStatus
}

type UpdateEnrollmentStatusUseCase struct {
	enrollmentRepo repositories.EnrollmentRepository
	publisher      messaging.Publisher
	logger         *logger.Logger
}

func NewUpdateEnrollmentStatusUseCase(
	enrollmentRepo repositories.EnrollmentRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *UpdateEnrollmentStatusUseCase {
	return &UpdateEnrollmentStatusUseCase{
		enrollmentRepo: enrollmentRepo,
		publisher:      publisher,
		logger:         logger,
	}
}

func (uc *UpdateEnrollmentStatusUseCase) Execute(ctx context.Context, input UpdateEnrollmentStatusInput) (*dtos.EnrollmentDTO, error) {
	enrollment, err := uc.enrollmentRepo.FindByID(ctx, input.EnrollmentID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, ErrEnrollmentNotFound
	}

	if !input.Status.IsValid() {
		return nil, valueobjects.ErrInvalidEnrollmentStatus
	}

	enrollment.Status = input.Status
	enrollment.UpdatedAt = time.Now().UTC()

	err = uc.enrollmentRepo.Update(ctx, enrollment)
	if err != nil {
		return nil, err
	}

	event := events.EnrollmentUpdatedEvent{
		ID:                 enrollment.ID,
		StudentID:          enrollment.StudentID,
		StudentUsername:    enrollment.StudentUsername,
		CourseID:           enrollment.CourseID,
		CourseName:         enrollment.CourseName,
		CourseOfferingID:   enrollment.CourseOfferingID,
		CourseOfferingName: enrollment.CourseOfferingName,
		Status:             enrollment.Status.String(),
		CreatedAt:          enrollment.CreatedAt,
		UpdatedAt:          enrollment.UpdatedAt,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeEnrollmentUpdated, event); err != nil {
			uc.logger.Error("Failed to publish enrollment updated event", zap.Error(err))
		}
	}

	uc.logger.Info("Enrollment Status Updated Successfully.",
		zap.String("enrollment_id", enrollment.ID),
		zap.String("status", enrollment.Status.String()),
	)

	var dto dtos.EnrollmentDTO
	dto.FromEntity(enrollment)
	return &dto, nil
}

