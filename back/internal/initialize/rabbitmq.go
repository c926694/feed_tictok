package initialize

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ(cfg RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	RabbitConn = conn
	RabbitChannel = ch
	return conn, ch, nil
}

func CloseRabbitMQ() {
	if RabbitChannel != nil {
		_ = RabbitChannel.Close()
	}
	if RabbitConn != nil {
		_ = RabbitConn.Close()
	}
}
