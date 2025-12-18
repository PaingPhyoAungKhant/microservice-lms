package mocks

import (
	"context"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"github.com/stretchr/testify/mock"
)

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) StoreAccessToken(ctx context.Context, userID string, accessToken string, expiration time.Duration) error {
	args := m.Called(ctx, userID, accessToken, expiration)
	return args.Error(0)
}

func (m *MockRedis) GetUserFromAccessToken(ctx context.Context, accessToken string) (string, error) {
	args := m.Called(ctx, accessToken)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) RevokeAccessToken(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

func (m *MockRedis) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiration time.Duration) error {
	args := m.Called(ctx, userID, refreshToken, expiration)
	return args.Error(0)
}

func (m *MockRedis) GetUserFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockRedis) StoreUserSession(ctx context.Context, sessionID, userID, accessToken, refreshToken, ipAddress, userAgent string, expiration time.Duration) error {
	args := m.Called(ctx, sessionID, userID, accessToken, refreshToken, ipAddress, userAgent, expiration)
	return args.Error(0)
}

func (m *MockRedis) GetUserSession(ctx context.Context, sessionID string) (*utils.SessionData, error) {
	args := m.Called(ctx, sessionID)
	session, _ := args.Get(0).(*utils.SessionData)
	return session, args.Error(1)
}

func (m *MockRedis) DeleteUserSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockRedis) ExpireUserSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	args := m.Called(ctx, sessionID, expiration)
	return args.Error(0)
}

func (m *MockRedis) UpdateLastActivity(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockRedis) StoreForgotPasswordOTP(ctx context.Context, userID, otp string) error {
	args := m.Called(ctx, userID, otp)
	return args.Error(0)
}

func (m *MockRedis) GetUserFromForgotPasswordOTP(ctx context.Context, otp string) (string, error) {
	args := m.Called(ctx, otp)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) RevokeForgotPasswordOTP(ctx context.Context, otp string) error {
	args := m.Called(ctx, otp)
	return args.Error(0)
}

func (m *MockRedis) StoreResetPasswordToken(ctx context.Context, userID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockRedis) GetUserFromResetPasswordToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) RevokeResetPasswordToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRedis) StoreVerifyEmailToken(ctx context.Context, userID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockRedis) GetUserFromVerifyEmailToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) RevokeVerifyEmailToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

