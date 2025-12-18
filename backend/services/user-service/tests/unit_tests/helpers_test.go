package unit_test

import (
	"testing"

	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/stretchr/testify/require"
)

func assertDTOEqualUser(t *testing.T, dto *dtos.UserDTO, user *entities.User) {
	t.Helper()

	require.NotNil(t, dto)
	require.Equal(t, user.ID, dto.ID)
	require.Equal(t, user.Email.String(), dto.Email)
	require.Equal(t, user.Username, dto.Username)
	require.Equal(t, user.Role.String(), dto.Role)
	require.Equal(t, user.Status.String(), dto.Status)
}

