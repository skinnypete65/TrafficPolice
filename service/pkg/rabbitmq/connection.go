package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func NewRabbitMQConn() (*amqp.Connection, error) {
	return amqp.Dial(os.Getenv("RABBITMQ_URL"))
}
