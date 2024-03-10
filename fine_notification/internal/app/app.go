package app

import (
	"fine_notification/internal/services"
	"fine_notification/internal/transport/rabbitmq"
	mqcommon "fine_notification/pkg/rabbitmq"
	"log"
)

func Run() {
	fineService := services.NewFineService()
	fineConsumer := setupFineConsumer(fineService)

	err := fineConsumer.StartConsume(mqcommon.ConsumeParams{
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

func setupFineConsumer(fineService services.FineService) *rabbitmq.FineConsumer {
	fineConsumer, err := rabbitmq.NewFineConsumer(fineService)
	if err != nil {
		log.Fatal(err)
	}
	err = fineConsumer.SetupExchangeAndQueue(
		mqcommon.ExchangeParams{
			Name:       mqcommon.FineExchange,
			Kind:       "fanout",
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		}, mqcommon.QueueParams{
			Name:       "fine_queue",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		mqcommon.BindingParams{
			Queue:    "fine_queue",
			Key:      "",
			Exchange: mqcommon.FineExchange,
			NoWait:   false,
			Args:     nil,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return fineConsumer
}
