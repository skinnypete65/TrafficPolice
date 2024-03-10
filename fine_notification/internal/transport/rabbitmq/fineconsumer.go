package rabbitmq

import (
	"TrafficPolice/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type FineConsumer struct {
	amqpChan *amqp.Channel
}

func NewFineConsumer() (*FineConsumer, error) {
	mqConn, err := rabbitmq.NewRabbitMQConn()
	if err != nil {
		return nil, err
	}

	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &FineConsumer{
		amqpChan: amqpChan,
	}, nil
}
