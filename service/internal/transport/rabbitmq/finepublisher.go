package rabbitmq

import (
	"TrafficPolice/internal/transport/rest/dto"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	jsonContentType = "application/json"
	FineExchange    = "fine"
	FineQueue       = "fine_queue"
	Fanout          = "fanout"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name FinePublisher
type FinePublisher interface {
	PublishFineNotification(c dto.CaseWithImage) error
}

type FinePublisherRabbitMQ struct {
	amqpChan *amqp.Channel
}

func NewFinePublisher(mqConn *amqp.Connection) (*FinePublisherRabbitMQ, error) {
	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &FinePublisherRabbitMQ{
		amqpChan: amqpChan,
	}, nil
}

func (p *FinePublisherRabbitMQ) SetupExchangeAndQueue(
	exchangeParams ExchangeParams,
	queueParams QueueParams,
	bindingsParams BindingParams,
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

	_, err = p.amqpChan.QueueDeclare(
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
		bindingsParams.Queue,
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

func (p *FinePublisherRabbitMQ) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		log.Printf("FinePublisherRabbitMQ CloseChan: %v\n", err)
	}
}

func (p *FinePublisherRabbitMQ) Publish(exchange string, contentType string, body []byte) error {
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

func (p *FinePublisherRabbitMQ) PublishFineNotification(c dto.CaseWithImage) error {
	cBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return p.Publish(FineExchange, jsonContentType, cBytes)
}
