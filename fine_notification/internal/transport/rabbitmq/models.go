package rabbitmq

import "github.com/rabbitmq/amqp091-go"

const (
	FineExchange = "fine"
)

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
	Queue    string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp091.Table
}

type ConsumeParams struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp091.Table
}
