package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type QueueParams struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp091.Table
}

type ExchangeParams struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp091.Table
}

type BindingParams struct {
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp091.Table
}
