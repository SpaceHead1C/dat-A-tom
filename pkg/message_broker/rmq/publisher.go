package rmq

import (
	"context"
	"errors"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type Publisher struct {
	*rmq.Publisher
}

type Publishing struct {
	p           Publisher
	routingKeys []string
	msg         []byte
	opts        []func(*rmq.PublishOptions)
}

type PublisherConfig struct {
	Logger   *zap.SugaredLogger
	Conn     *Connection
	Exchange string
}

func NewPublisher(c PublisherConfig) (*Publisher, error) {
	if c.Conn == nil {
		return nil, errors.New("broker connection can not be nil")
	}
	p, err := rmq.NewPublisher(
		c.Conn.conn,
		rmq.WithPublisherOptionsExchangeName(c.Exchange),
		rmq.WithPublisherOptionsLogger(logger{c.Logger}),
	)
	return &Publisher{
		Publisher: p,
	}, err
}

func (p *Publisher) Close() {
	p.Publisher.Close()
}

func NewPublishing(p *Publisher, routingKeys []string, msg []byte, opts ...func(*rmq.PublishOptions)) *Publishing {
	return &Publishing{
		p:           *p,
		routingKeys: routingKeys,
		msg:         msg,
		opts:        opts,
	}
}

func (p *Publishing) Publish(ctx context.Context) error {
	return p.p.PublishWithContext(
		ctx,
		p.msg,
		p.routingKeys,
		p.opts...,
	)
}
