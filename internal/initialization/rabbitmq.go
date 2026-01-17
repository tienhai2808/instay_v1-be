package initialization

import (
	"fmt"

	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	Conn *amqp091.Connection
	Chan *amqp091.Channel
}

func InitRabbitMQ(cfg *config.Config) (*MQ, error) {
	protocol := "amqp"
	if cfg.RabbitMQ.UseSSL {
		protocol = protocol + "s"
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		protocol,
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Vhost,
	)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("message queue - %w", err)
	}

	chann, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("message queue - %w", err)
	}

	return &MQ{
		conn,
		chann,
	}, nil
}

func (mq *MQ) Close() {
	_ = mq.Chan.Close()
	_ = mq.Conn.Close()
}
