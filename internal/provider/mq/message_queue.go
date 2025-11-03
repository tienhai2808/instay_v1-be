package mq

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MessageQueueProvider interface {
	PublishMessage(exchange, routingKey string, body []byte) error
	ConsumeMessage(queueName, exchange, routingKey string, handler func([]byte) error) error
}

type messageQueueProviderImpl struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	logger *zap.Logger
}

func NewMessageQueueProvider(
	conn *amqp091.Connection,
	ch *amqp091.Channel,
	logger *zap.Logger,
) MessageQueueProvider {
	return &messageQueueProviderImpl{
		conn,
		ch,
		logger,
	}
}

func (m *messageQueueProviderImpl) PublishMessage(exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         body,
	}); err != nil {
		return err
	}

	return nil
}

func (m *messageQueueProviderImpl) ConsumeMessage(queueName, exchange, routingKey string, handler func([]byte) error) error {
	if _, err := m.ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}

	if err := m.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	if err := m.ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return err
	}

	if err := m.ch.Qos(5, 0, false); err != nil {
		return err
	}

	msgs, err := m.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for i := range 5 {
		go func(workerID int) {
			for msg := range msgs {
				m.processWithRetry(msg.Body, handler, workerID)
			}
		}(i)
	}

	return nil
}

func (m *messageQueueProviderImpl) processWithRetry(body []byte, handler func([]byte) error, workerID int) {
	maxAttempts := 5
	initialInterval := 1000 * time.Millisecond
	multiplier := 2.0
	maxInterval := 10000 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := handler(body)
		if err == nil {
			return
		}
		m.logger.Error(fmt.Sprintf("work %d (%d/%d) failed", workerID, attempt, maxAttempts), zap.Error(err))

		if attempt < maxAttempts {
			delay := float64(initialInterval) * math.Pow(multiplier, float64(attempt-1))
			if delay > float64(maxInterval) {
				delay = float64(maxInterval)
			}
			time.Sleep(time.Duration(delay))
		}
	}

	m.logger.Error(fmt.Sprintf("work %d", workerID), zap.Error(fmt.Errorf("message sending failed after %d attempts", maxAttempts)))
}
