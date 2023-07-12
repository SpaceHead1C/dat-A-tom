package amq

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

const QueueDLEArg = "x-dead-letter-exchange"

type ConsumerConfig struct {
	Logger    *zap.SugaredLogger
	Conn      *rmq.Conn
	Queue     string
	Handler   rmq.Handler
	QueueArgs QueueArgs
}

func RunNewConsumer(c ConsumerConfig) error {
	consumer, err := rmq.NewConsumer(
		c.Conn,
		c.Handler,
		c.Queue,
		rmq.WithConsumerOptionsQueueDurable,
		rmq.WithConsumerOptionsQueueArgs(c.QueueArgs.asRmqTable()),
		rmq.WithConsumerOptionsLogger(logger{c.Logger}),
	)
	if err != nil {
		return err
	}
	defer consumer.Close()
	var forever chan struct{}
	<-forever
	return errors.New("rmq listener is down")
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

func (qa QueueArgs) asRmqTable() rmq.Table {
	return rmq.Table(qa)
}
