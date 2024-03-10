package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConn() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
}
