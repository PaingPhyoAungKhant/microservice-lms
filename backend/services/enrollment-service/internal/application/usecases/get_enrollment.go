package usecases

import (
	"context"
	"errors"

	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
)

var (
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrInvalidEnrollmentID = errors.New("invalid enrollment ID")
)

type GetEnrollmentInput struct {
	EnrollmentID string
}

type GetEnrollmentUseCase struct {
	enrollmentRepo repositories.EnrollmentRepository
	logger         *logger.Logger
}

func NewGetEnrollmentUseCase(
	enrollmentRepo repositories.EnrollmentRepository,
	logger *logger.Logger,
) *GetEnrollmentUseCase {
	return &GetEnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (uc *GetEnrollmentUseCase) Execute(ctx context.Context, input GetEnrollmentInput) (*dtos.EnrollmentDTO, error) {
	enrollment, err := uc.enrollmentRepo.FindByID(ctx, input.EnrollmentID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, ErrEnrollmentNotFound
	}

	var dto dtos.EnrollmentDTO
	dto.FromEntity(enrollment)
	return &dto, nil
}

