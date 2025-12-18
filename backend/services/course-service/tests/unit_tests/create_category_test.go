package unit_test

import (
	"context"
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/tests/mocks"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory_CategoryAlreadyExists(t *testing.T) {
	existing := entities.NewCategory("Existing Category", "Description")
	repo := new(mocks.MockCategoryRepository)
	repo.On("FindByName", mock.Anything, "Existing Category").Return(existing, nil).Once()
	logger := logger.NewNop()
	uc := usecases.NewCreateCategoryUseCase(repo, logger)

	_, err := uc.Execute(context.Background(), dtos.CreateCategoryInput{
		Name:        "Existing Category",
		Description: "New Description",
	})
	require.ErrorIs(t, err, usecases.ErrCategoryAlreadyExists)

	repo.AssertExpectations(t)
}

func TestCreateCategory_Success(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()
	input := dtos.CreateCategoryInput{
		Name:        "New Category",
		Description: "Category Description",
	}

	var createdCategory *entities.Category
	repo.On("FindByName", mock.Anything, input.Name).Return((*entities.Category)(nil), nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.Category) bool {
		createdCategory = c
		return true
	})).Return(nil).Once()

	uc := usecases.NewCreateCategoryUseCase(repo, logger)

	dto, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, createdCategory)
	require.Equal(t, input.Name, createdCategory.Name)
	require.Equal(t, input.Description, createdCategory.Description)

	assertDTOEqualCategory(t, dto, createdCategory)
	repo.AssertExpectations(t)
}

func TestCreateCategory_EmptyName(t *testing.T) {
	repo := new(mocks.MockCategoryRepository)
	logger := logger.NewNop()
	uc := usecases.NewCreateCategoryUseCase(repo, logger)

	_, err := uc.Execute(context.Background(), dtos.CreateCategoryInput{
		Name:        "",
		Description: "Description",
	})
	require.Error(t, err)

	repo.AssertNotCalled(t, "FindByName", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

