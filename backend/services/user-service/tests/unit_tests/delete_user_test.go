package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser_FindByIDError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	repo.On("FindByID", mock.Anything, "user-1").Return((*entities.User)(nil), expectedErr).Once()

	uc := usecases.NewDeleteUserUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteUserInput{UserID: "user-1"})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteUser_AdminForbidden(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("admin@example.com")
	roleVO, _ := valueobjects.NewRole("admin")
	user := entities.NewUser(emailVO, "admin", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()

	uc := usecases.NewDeleteUserUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.DeleteUserInput{UserID: user.ID})
	require.ErrorIs(t, err, usecases.ErrCannotDeleteAdminUser)

	publisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	emailVO, _ := valueobjects.NewEmail("user@example.com")
	roleVO, _ := valueobjects.NewRole("student")
	user := entities.NewUser(emailVO, "student", roleVO, "hash")

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)

	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	repo.On("Delete", mock.Anything, user.ID).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeUserDeleted, mock.Anything).Return(nil).Once()

	uc := usecases.NewDeleteUserUseCase(repo, publisher, logger.NewNop())

	out, err := uc.Execute(context.Background(), usecases.DeleteUserInput{UserID: user.ID})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotEmpty(t, out.Message)

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.UserDeletedEvent)
	require.True(t, ok)
	assert.Equal(t, user.ID, event.ID)
	assert.WithinDuration(t, time.Now(), event.DeletedAt, 2*time.Second)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

