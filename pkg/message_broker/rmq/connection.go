package rmq

import (
	"datatom/pkg/log"
	"errors"
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

type Connection struct {
	conn *rmq.Conn
	l    *zap.SugaredLogger
}

func NewConnection(c ConnectionConfig) (*Connection, error) {
	if strings.TrimSpace(c.Address) == "" {
		return nil, errors.New("RMQ address can not be empty")
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
	conn, err := rmq.NewConn(
		connectionString(c),
		rmq.WithConnectionOptionsLogger(logger{c.Logger}),
	)
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn: conn,
		l:    c.Logger,
	}, nil
}

func (conn *Connection) Close() {
	if err := conn.conn.Close(); err != nil {
		conn.l.Errorf("rmq connection close error: %s", err)
	}
}

func connectionString(c ConnectionConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.User, c.Password, c.Address, c.Port, c.VHost)
}
