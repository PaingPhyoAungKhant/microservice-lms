package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, routingKey string, message interface{}) error {
	args := m.Called(ctx, routingKey, message)
	return args.Error(0)
}

