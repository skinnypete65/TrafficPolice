package rabbitmq

import (
	"TrafficPolice/pkg/rabbitmq"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type FinePublisher struct {
	amqpChan *amqp.Channel
}

func NewFinePublisher() (*FinePublisher, error) {
	mqConn, err := rabbitmq.NewRabbitMQConn()
	if err != nil {
		return nil, err
	}

	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &FinePublisher{
		amqpChan: amqpChan,
	}, nil
}

func (p *FinePublisher) SetupExchangeAndQueue(
	exchangeParams rabbitmq.ExchangeParams,
	queueParams rabbitmq.QueueParams,
	bindingsParams rabbitmq.BindingParams,
) error {
	err := p.amqpChan.ExchangeDeclare(
		exchangeParams.Name,
		exchangeParams.Kind,
		exchangeParams.Durable,
		exchangeParams.AutoDelete,
		exchangeParams.Internal,
		exchangeParams.NoWait,
		exchangeParams.Args,
	)

	if err != nil {
		return err
	}

	queue, err := p.amqpChan.QueueDeclare(
		queueParams.Name,
		queueParams.Durable,
		queueParams.AutoDelete,
		queueParams.Exclusive,
		queueParams.NoWait,
		queueParams.Args,
	)
	if err != nil {
		return err
	}

	err = p.amqpChan.QueueBind(
		queue.Name,
		bindingsParams.Key,
		bindingsParams.Exchange,
		bindingsParams.NoWait,
		bindingsParams.Args,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *FinePublisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		log.Printf("EmailsPublisher CloseChan: %v\n", err)
	}
}

func (p *FinePublisher) Publish(exchange string, contentType string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.amqpChan.PublishWithContext(
		ctx,
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)

	return err
}
