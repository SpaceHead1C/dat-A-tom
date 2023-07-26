package rmq

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

const QueueDLEArg = "x-dead-letter-exchange"

type ConsumerConfig struct {
	Logger    *zap.SugaredLogger
	Conn      *Connection
	Queue     string
	Handler   rmq.Handler
	QueueArgs QueueArgs
}

type Consumer struct {
	conn    *Connection
	queue   string
	handler rmq.Handler
	args    QueueArgs
	l       *zap.SugaredLogger
}

func NewConsumer(c ConsumerConfig) *Consumer {
	return &Consumer{
		conn:    c.Conn,
		queue:   c.Queue,
		handler: c.Handler,
		args:    c.QueueArgs,
		l:       c.Logger,
	}
}

func (c *Consumer) Consume() error {
	consumer, err := rmq.NewConsumer(
		c.conn.conn,
		c.handler,
		c.queue,
		rmq.WithConsumerOptionsQueueDurable,
		rmq.WithConsumerOptionsQueueArgs(c.args.asRmqTable()),
		rmq.WithConsumerOptionsLogger(logger{c.l}),
	)
	if err != nil {
		return err
	}
	defer consumer.Close()
	var forever chan struct{}
	<-forever
	return errors.New("rmq consumer is down")
}

type QueueArgs rmq.Table

func NewQueueArgs() QueueArgs {
	return make(QueueArgs)
}

func (qa QueueArgs) AddArg(key string, value any) QueueArgs {
	qa[key] = value
	return qa
}

func (qa QueueArgs) AddDLEArg(value string) QueueArgs {
	return qa.AddArg(QueueDLEArg, value)
}

func (qa QueueArgs) AddTypeArg(value string) QueueArgs {
	return qa.AddArg(amqp.QueueTypeArg, value)
}

func (qa QueueArgs) AsClassic() QueueArgs {
	return qa.AddTypeArg(amqp.QueueTypeClassic)
}

func (qa QueueArgs) asRmqTable() rmq.Table {
	return rmq.Table(qa)
}
