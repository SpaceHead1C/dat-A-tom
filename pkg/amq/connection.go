package amq

import (
	"datatom/pkg/log"
	"fmt"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"strings"
)

const (
	defaultUserRMQ     = "guest"
	defaultPasswordRMQ = "guest"
)

type ConnectionConfig struct {
	Logger   *zap.SugaredLogger
	Address  string
	Port     uint
	User     string
	Password string
	VHost    string
}

func NewConnection(c ConnectionConfig) (*rmq.Conn, error) {
	if strings.TrimSpace(c.Address) == "" {
		return nil, nil
	}
	if c.Port == 0 {
		c.Port = 5672
	}
	if c.Logger == nil {
		c.Logger = log.GlobalLogger()
	}
	if strings.TrimSpace(c.User) == "" {
		c.User = defaultUserRMQ
	}
	if strings.TrimSpace(c.Password) == "" {
		c.Password = defaultPasswordRMQ
	}
	return rmq.NewConn(
		connectionString(c),
		rmq.WithConnectionOptionsLogger(logger{c.Logger}),
	)
}

func connectionString(c ConnectionConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.User, c.Password, c.Address, c.Port, c.VHost)
}
