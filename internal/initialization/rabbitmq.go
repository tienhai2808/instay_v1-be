package initialization

import (
	"fmt"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	Conn *amqp091.Connection
	Chan *amqp091.Channel
}

func InitRabbitMQ(cfg *config.Config) (*MQ, error) {
	dsn := fmt.Sprintf("amqps://%s:%s@%s/%s",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Vhost,
	)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, err
	}

	chann, err := conn.Channel()
	if err != nil {
		return nil, err
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
