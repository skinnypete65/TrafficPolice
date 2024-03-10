package rabbitmq

import (
	"fine_notification/internal/services"
	"fine_notification/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type FineConsumer struct {
	amqpChan    *amqp.Channel
	fineService services.FineService
}

func NewFineConsumer(fineService services.FineService) (*FineConsumer, error) {
	mqConn, err := rabbitmq.NewRabbitMQConn()
	if err != nil {
		return nil, err
	}

	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &FineConsumer{
		amqpChan:    amqpChan,
		fineService: fineService,
	}, nil
}

func (p *FineConsumer) SetupExchangeAndQueue(
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

func (p *FineConsumer) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		log.Printf("EmailsPublisher CloseChan: %v\n", err)
	}
}

func (p *FineConsumer) StartConsume(params rabbitmq.ConsumeParams) error {
	msgs, err := p.amqpChan.Consume(
		params.Queue,     // queue
		params.Consumer,  // consumer
		params.AutoAck,   // auto-ack
		params.Exclusive, // exclusive
		params.NoLocal,   // no-local
		params.NoWait,    // no-wait
		params.Args,      // args
	)

	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s\n", d.Body)
		}
	}()

	<-forever
	return err
}
