package messaging

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, routingKey string, message interface{}) error
}
