package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnrollment_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	logger := logger.NewNop()

	ctx := context.Background()

	enrollment := entities.NewEnrollment(
		uuid.New().String(),
		"getstudent",
		uuid.New().String(),
		"Get Course",
		uuid.New().String(),
		"Fall 2024",
	)
	enrollmentRepo.Create(ctx, enrollment)

	getUC := usecases.NewGetEnrollmentUseCase(enrollmentRepo, logger)

	input := usecases.GetEnrollmentInput{
		EnrollmentID: enrollment.ID,
	}

	result, err := getUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, enrollment.ID, result.ID)
	assert.Equal(t, enrollment.StudentID, result.StudentID)
	assert.Equal(t, enrollment.StudentUsername, result.StudentUsername)
}

func TestGetEnrollment_Integration_NotFound(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	logger := logger.NewNop()

	getUC := usecases.NewGetEnrollmentUseCase(enrollmentRepo, logger)

	input := usecases.GetEnrollmentInput{
		EnrollmentID: "00000000-0000-0000-0000-000000000000",
	}

	_, err := getUC.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, usecases.ErrEnrollmentNotFound, err)
}

