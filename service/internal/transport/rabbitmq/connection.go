package rabbitmq

import (
	"TrafficPolice/internal/config"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConn(cfg *config.Config) (*amqp.Connection, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	return amqp.Dial(url)
}
