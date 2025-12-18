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

func TestFindEnrollment_Integration_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	logger := logger.NewNop()

	ctx := context.Background()

	enrollment1 := entities.NewEnrollment(
		uuid.New().String(),
		"findstudent1",
		uuid.New().String(),
		"Find Course 1",
		uuid.New().String(),
		"Fall 2024",
	)
	enrollment2 := entities.NewEnrollment(
		uuid.New().String(),
		"findstudent2",
		uuid.New().String(),
		"Find Course 2",
		uuid.New().String(),
		"Spring 2024",
	)
	enrollmentRepo.Create(ctx, enrollment1)
	enrollmentRepo.Create(ctx, enrollment2)

	findUC := usecases.NewFindEnrollmentUseCase(enrollmentRepo, logger)

	limit := 10
	input := usecases.FindEnrollmentInput{
		Limit: &limit,
	}

	result, err := findUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 2)
	assert.GreaterOrEqual(t, len(result.Enrollments), 2)
}

func TestFindEnrollment_Integration_WithFilters(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	enrollmentRepo := SetupEnrollmentRepository(db)
	logger := logger.NewNop()

	ctx := context.Background()

	enrollment := entities.NewEnrollment(
		uuid.New().String(),
		"filterspecificstudent",
		uuid.New().String(),
		"Filter Specific Course",
		uuid.New().String(),
		"Fall 2024",
	)
	enrollmentRepo.Create(ctx, enrollment)

	findUC := usecases.NewFindEnrollmentUseCase(enrollmentRepo, logger)

	limit := 10
	input := usecases.FindEnrollmentInput{
		StudentID: &enrollment.StudentID,
		Limit:     &limit,
	}

	result, err := findUC.Execute(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 1)
}

