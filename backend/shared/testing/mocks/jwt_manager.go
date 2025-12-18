package mocks

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/mock"
)

type MockJwtManager struct {
	mock.Mock
}

func (m *MockJwtManager) GenerateAccessToken(userId, email, role, status string) (string, error) {
	args := m.Called(userId, email, role, status)
	return args.String(0), args.Error(1)
}

func (m *MockJwtManager) GenerateRefreshToken(userId, email, role, status string) (string, error) {
	args := m.Called(userId, email, role, status)
	return args.String(0), args.Error(1)
}

func (m *MockJwtManager) VerifyToken(tokenString string) (utils.JwtClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(utils.JwtClaims), args.Error(1)
}

func (m *MockJwtManager) AccessTokenDuration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockJwtManager) RefreshTokenDuration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

