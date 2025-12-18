package consumer

import (
	"context"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

type EventConsumer struct {
	rabbitMQ              *messaging.RabbitMQ
	userUpdatedHandler    *handlers.UserUpdatedHandler
	zoomMeetingCreatedHandler *handlers.ZoomMeetingCreatedHandler
	logger                *logger.Logger
}

func NewEventConsumer(
	rabbitMQ *messaging.RabbitMQ,
	userUpdatedHandler *handlers.UserUpdatedHandler,
	zoomMeetingCreatedHandler *handlers.ZoomMeetingCreatedHandler,
	logger *logger.Logger,
) *EventConsumer {
	return &EventConsumer{
		rabbitMQ:              rabbitMQ,
		userUpdatedHandler:    userUpdatedHandler,
		zoomMeetingCreatedHandler: zoomMeetingCreatedHandler,
		logger:                logger,
	}
}

func (c *EventConsumer) Start(ctx context.Context) error {
	routingKeys := []string{
		events.EventTypeUserUpdated,
		events.EventTypeZoomMeetingCreated,
	}

	messages, err := c.rabbitMQ.Consume(ctx, "course-service.queue", routingKeys)
	if err != nil {
		return err
	}

	c.logger.Info("started consuming events",
		zap.Strings("routing_keys", routingKeys),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopping event consumer")
			return nil
		case msg, ok := <-messages:
			if !ok {
				c.logger.Warn("message channel closed")
				return nil
			}

			if err := c.handleMessage(msg); err != nil {
				c.logger.Error("failed to handle message",
					zap.String("routing_key", msg.RoutingKey),
					zap.Error(err),
				)
				c.rabbitMQ.Reject(msg.Delivery, true)
			} else {
				c.rabbitMQ.Acknowledge(msg.Delivery)
			}
		}
	}
}

func (c *EventConsumer) handleMessage(msg messaging.Message) error {
	c.logger.Info("received message",
		zap.String("routing_key", msg.RoutingKey),
	)

	switch msg.RoutingKey {
	case events.EventTypeUserUpdated:
		return c.userUpdatedHandler.Handle(msg.Body)
	case events.EventTypeZoomMeetingCreated:
		return c.zoomMeetingCreatedHandler.Handle(msg.Body)
	default:
		c.logger.Warn("unknown routing key",
			zap.String("routing_key", msg.RoutingKey),
		)
		return nil
	}
}

