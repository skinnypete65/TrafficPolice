package app

import (
	"fine_notification/internal/config"
	"fine_notification/internal/mailer"
	"fine_notification/internal/transport/rabbitmq"
	"gopkg.in/gomail.v2"
	"log"
)

func Run() {
	cfg, err := config.ParseConfig("notification_config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	dialer := gomail.NewDialer("smtp.gmail.com", 587, cfg.EmailSenderUsername, cfg.EmailSenderPass)

	fineMailer := mailer.NewMailer(dialer)
	fineConsumer := setupFineConsumer(fineMailer)

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

func setupFineConsumer(mailer *mailer.Mailer) *rabbitmq.FineConsumer {
	fineConsumer, err := rabbitmq.NewFineConsumer(mailer)
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
