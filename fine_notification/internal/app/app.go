package app

import (
	"fine_notification/internal/config"
	"fine_notification/internal/mailer"
	"fine_notification/internal/transport/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
	"log"
)

const (
	notificationConfigPath = "notification_config.yaml"
)

func Run() {
	cfg, err := config.ParseConfig(notificationConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	mqConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatal(err)
	}

	dialer := setupMailDialer(cfg)

	fineMailer := mailer.NewMailer(dialer)
	fineConsumer := setupFineConsumer(mqConn, fineMailer)

	err = fineConsumer.StartConsume(rabbitmq.ConsumeParams{
		Queue:     "fine_queue",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func setupFineConsumer(mqConn *amqp.Connection, mailer *mailer.Mailer) *rabbitmq.FineConsumer {
	fineConsumer, err := rabbitmq.NewFineConsumer(mqConn, mailer)
	if err != nil {
		log.Fatal(err)
	}
	err = fineConsumer.SetupExchangeAndQueue(
		rabbitmq.ExchangeParams{
			Name:       rabbitmq.FineExchange,
			Kind:       "fanout",
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		}, rabbitmq.QueueParams{
			Name:       "fine_queue",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		rabbitmq.BindingParams{
			Queue:    "fine_queue",
			Key:      "",
			Exchange: rabbitmq.FineExchange,
			NoWait:   false,
			Args:     nil,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return fineConsumer
}

func setupMailDialer(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(
		cfg.EmailSender.Host,
		cfg.EmailSender.Port,
		cfg.EmailSender.Username,
		cfg.EmailSender.Password,
	)
}
