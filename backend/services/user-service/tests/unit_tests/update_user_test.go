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

func TestUpdateUser_RepoError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := new(mocks.MockUserRepository)
	repo.On("FindByID", mock.Anything, "user-123").Return((*entities.User)(nil), expectedErr).Once()
	publisher := new(mocks.MockPublisher)

	uc := usecases.NewUpdateUserUseCase(repo, publisher, logger.NewNop())
	_, err := uc.Execute(context.Background(), usecases.UpdateUserInput{ID: "user-123"})
	require.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	user := newTestUser(t, "user@example.com", "student")
	repo := new(mocks.MockUserRepository)
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	publisher := new(mocks.MockPublisher)

	uc := usecases.NewUpdateUserUseCase(repo, publisher, logger.NewNop())
	newEmail := "bad-email"
	_, err := uc.Execute(context.Background(), usecases.UpdateUserInput{
		ID:    user.ID,
		Email: &newEmail,
	})
	require.ErrorIs(t, err, valueobjects.ErrInvalidEmail)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateUser_InvalidRole(t *testing.T) {
	user := newTestUser(t, "user@example.com", "student")
	repo := new(mocks.MockUserRepository)
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	publisher := new(mocks.MockPublisher)

	uc := usecases.NewUpdateUserUseCase(repo, publisher, logger.NewNop())
	newRole := "invalid"
	_, err := uc.Execute(context.Background(), usecases.UpdateUserInput{
		ID:   user.ID,
		Role: &newRole,
	})
	require.ErrorIs(t, err, valueobjects.ErrInvalidRole)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateUser_InvalidStatus(t *testing.T) {
	user := newTestUser(t, "user@example.com", "student")
	repo := new(mocks.MockUserRepository)
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	publisher := new(mocks.MockPublisher)

	uc := usecases.NewUpdateUserUseCase(repo, publisher, logger.NewNop())
	newStatus := "not-a-status"
	_, err := uc.Execute(context.Background(), usecases.UpdateUserInput{
		ID:     user.ID,
		Status: &newStatus,
	})
	require.ErrorIs(t, err, valueobjects.ErrInvalidStatus)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	now := time.Now().Add(-time.Hour)
	user := newTestUser(t, "user@example.com", "student")
	user.CreatedAt = now
	user.UpdatedAt = now

	repo := new(mocks.MockUserRepository)
	publisher := new(mocks.MockPublisher)

	var updatedUser *entities.User
	repo.On("FindByID", mock.Anything, user.ID).Return(user, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
		updatedUser = u
		return true
	})).Return(nil).Once()
	publisher.On("Publish", mock.Anything, events.EventTypeUserUpdated, mock.Anything).Return(nil).Once()

	newUsername := "updated"
	newEmail := "updated@example.com"
	newRole := "instructor"
	newStatus := "inactive"

	uc := usecases.NewUpdateUserUseCase(repo, publisher, logger.NewNop())
	dto, err := uc.Execute(context.Background(), usecases.UpdateUserInput{
		ID:       user.ID,
		Username: &newUsername,
		Email:    &newEmail,
		Role:     &newRole,
		Status:   &newStatus,
	})
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	assert.Equal(t, newUsername, updatedUser.Username)
	assert.Equal(t, newEmail, updatedUser.Email.String())
	assert.Equal(t, newRole, updatedUser.Role.String())
	assert.Equal(t, newStatus, updatedUser.Status.String())
	assert.True(t, updatedUser.UpdatedAt.After(now))

	require.Len(t, publisher.Calls, 1)
	event, ok := publisher.Calls[0].Arguments.Get(2).(events.UserUpdatedEvent)
	require.True(t, ok)
	assert.Equal(t, updatedUser.ID, event.ID)

	assertDTOEqualUser(t, dto, updatedUser)

	repo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}


func newTestUser(t *testing.T, email, role string) *entities.User {
	t.Helper()
	emailVO, err := valueobjects.NewEmail(email)
	require.NoError(t, err)
	roleVO, err := valueobjects.NewRole(role)
	require.NoError(t, err)
	return entities.NewUser(emailVO, "username", roleVO, "hash")
}

