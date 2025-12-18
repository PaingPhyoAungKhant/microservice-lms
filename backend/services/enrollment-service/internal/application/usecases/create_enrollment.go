package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

var (
	ErrEnrollmentAlreadyExists = errors.New("enrollment already exists")
)

type CreateEnrollmentInput struct {
	StudentID          string
	StudentUsername    string
	CourseID           string
	CourseName         string
	CourseOfferingID   string
	CourseOfferingName string
}

type CreateEnrollmentUseCase struct {
	enrollmentRepo repositories.EnrollmentRepository
	publisher      messaging.Publisher
	logger         *logger.Logger
}

func NewCreateEnrollmentUseCase(
	enrollmentRepo repositories.EnrollmentRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
) *CreateEnrollmentUseCase {
	return &CreateEnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		publisher:      publisher,
		logger:         logger,
	}
}

func (uc *CreateEnrollmentUseCase) Execute(ctx context.Context, input CreateEnrollmentInput) (*dtos.EnrollmentDTO, error) {
	limit := 1
	existing, _ := uc.enrollmentRepo.Find(ctx, repositories.EnrollmentQuery{
		StudentID:        &input.StudentID,
		CourseOfferingID: &input.CourseOfferingID,
		Limit:            &limit,
	})
	if existing != nil && len(existing.Enrollments) > 0 {
		return nil, ErrEnrollmentAlreadyExists
	}

	enrollment := entities.NewEnrollment(
		input.StudentID,
		input.StudentUsername,
		input.CourseID,
		input.CourseName,
		input.CourseOfferingID,
		input.CourseOfferingName,
	)

	if err := uc.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	event := events.EnrollmentCreatedEvent{
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
		if err := uc.publisher.Publish(ctx, events.EventTypeEnrollmentCreated, event); err != nil {
			uc.logger.Error("Failed to publish enrollment created event", zap.Error(err))
		}
	}

	uc.logger.Info("Enrollment Created Successfully.",
		zap.String("enrollment_id", enrollment.ID),
		zap.String("student_id", enrollment.StudentID),
		zap.String("course_offering_id", enrollment.CourseOfferingID),
	)

	var dto dtos.EnrollmentDTO
	dto.FromEntity(enrollment)
	return &dto, nil
}

