package usecases

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

type FindEnrollmentInput struct {
	SearchQuery      *string
	StudentID        *string
	CourseID         *string
	CourseOfferingID *string
	Status           *valueobjects.EnrollmentStatus
	Limit            *int
	Offset           *int
	SortColumn       *string
	SortDirection    *repositories.SortDirection
}

type FindEnrollmentOutput struct {
	Enrollments []dtos.EnrollmentDTO `json:"enrollments"`
	Total       int                   `json:"total"`
}

type FindEnrollmentUseCase struct {
	enrollmentRepo repositories.EnrollmentRepository
	logger         *logger.Logger
}

func NewFindEnrollmentUseCase(
	enrollmentRepo repositories.EnrollmentRepository,
	logger *logger.Logger,
) *FindEnrollmentUseCase {
	return &FindEnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (uc *FindEnrollmentUseCase) Execute(ctx context.Context, input FindEnrollmentInput) (*FindEnrollmentOutput, error) {
	query := repositories.EnrollmentQuery{
		SearchQuery:      input.SearchQuery,
		StudentID:        input.StudentID,
		CourseID:         input.CourseID,
		CourseOfferingID: input.CourseOfferingID,
		Status:           input.Status,
		Limit:            input.Limit,
		Offset:           input.Offset,
		SortColumn:       input.SortColumn,
		SortDirection:    input.SortDirection,
	}

	result, err := uc.enrollmentRepo.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	enrollmentDTOs := make([]dtos.EnrollmentDTO, len(result.Enrollments))
	for i, enrollment := range result.Enrollments {
		enrollmentDTOs[i].FromEntity(enrollment)
	}

	return &FindEnrollmentOutput{
		Enrollments: enrollmentDTOs,
		Total:      result.Total,
	}, nil
}

