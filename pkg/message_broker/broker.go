package message_broker

import "context"

type Connection interface {
	Close()
}

type Publisher interface {
	Close()
}

type Publishing interface {
	Publish(context.Context) error
}

type Consumer interface {
	Consume() error
}
