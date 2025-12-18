// Package messaging - shared package for RabbitMQ
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *logger.Logger
	config *config.RabbitMQConfig
}

func NewRabbitMQ(config *config.RabbitMQConfig, log *logger.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		config.Exchange,
		config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		logger:  log,
		config: config,
	}, nil
}

func (r *RabbitMQ) GetExchangeType() string {
	return r.config.ExchangeType
}

func (r *RabbitMQ) GetExchange() string {
	return r.config.Exchange
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	err = r.channel.PublishWithContext(
		ctx,
		r.config.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	r.logger.Info("message published",
		zap.String("routing_key", routingKey),
		zap.String("exchange", r.config.Exchange),
	)

	return nil
}

func (r RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}

	return nil
}

type Message struct {
	RoutingKey string
	Body       []byte
	Delivery   amqp.Delivery
}

func (r *RabbitMQ) Consume(ctx context.Context, queueName string, routingKeys []string) (<-chan Message, error) {
	queue, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, routingKey := range routingKeys {
		err = r.channel.QueueBind(
			queue.Name,
			routingKey,
			r.config.Exchange,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
		}
	}

	deliveries, err := r.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	messages := make(chan Message)
	go func() {
		defer close(messages)
		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-deliveries:
				if !ok {
					return
				}
				messages <- Message{	
					RoutingKey: delivery.RoutingKey,
					Body:       delivery.Body,
					Delivery:   delivery,
				}
			}
		}
	}()

	return messages, nil
}

func (r *RabbitMQ) Acknowledge(delivery amqp.Delivery) error {
	return delivery.Ack(false)
}

func (r *RabbitMQ) Reject(delivery amqp.Delivery, requeue bool) error {
	return delivery.Nack(false, requeue)
}
