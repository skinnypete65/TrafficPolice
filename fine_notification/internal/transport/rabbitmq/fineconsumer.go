package rabbitmq

import (
	"encoding/json"
	"fine_notification/internal/mailer"
	"fine_notification/internal/transport/dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type FineConsumer struct {
	amqpChan *amqp.Channel
	mailer   *mailer.Mailer
}

func NewFineConsumer(
	mqConn *amqp.Connection,
	mailer *mailer.Mailer,
) (*FineConsumer, error) {

	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &FineConsumer{
		amqpChan: amqpChan,
		mailer:   mailer,
	}, nil
}

func (p *FineConsumer) SetupExchangeAndQueue(
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

func (p *FineConsumer) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		log.Printf("EmailsPublisher CloseChan: %v\n", err)
	}
}

func (p *FineConsumer) StartConsume(params ConsumeParams) error {
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
			var cDto dto.Case
			err = json.Unmarshal(d.Body, &cDto)
			if err != nil {
				log.Println(err)
				continue
			}

			email := mailer.Email{
				From:    "vasyagoose8@gmail.com",
				To:      cDto.Transport.Person.Email,
				Subject: "Информация о совершенном правонарушении",
			}

			err = p.mailer.SendFineMessage(email, cDto)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	<-forever
	return err
}
