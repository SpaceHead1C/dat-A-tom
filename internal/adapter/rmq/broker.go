package rmq

import (
	"datatom/pkg/message_broker/rmq"
	"go.uber.org/zap"
)

type Config struct {
	Publisher *rmq.Publisher
	Logger    *zap.SugaredLogger
}

type Broker struct {
	publisher *rmq.Publisher
	l         *zap.SugaredLogger
}

func NewBroker(c Config) (*Broker, error) {
	return &Broker{
		publisher: c.Publisher,
		l:         c.Logger,
	}, nil
}
