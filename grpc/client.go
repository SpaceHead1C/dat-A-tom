package grpc

import (
	"datatom/internal/pb"
	"datatom/pkg/log"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Cli    pb.DatawayClient
	conn   *grpc.ClientConn
	logger *zap.SugaredLogger
}

type Connection struct {
	address string
	logger  *zap.SugaredLogger
}

type Config struct {
	Address string
	Port    uint
	Logger  *zap.SugaredLogger
}

func NewConnection(c Config) *Connection {
	if c.Address == "" {
		return nil
	}
	var err error
	l := c.Logger
	if l == nil {
		l, err = log.NewLogger()
		if err != nil {
			l = zap.L().Sugar()
		}
	}
	return &Connection{
		address: fmt.Sprintf("%s:%d", c.Address, c.Port),
		logger:  l,
	}
}

func (c Connection) dial() (*grpc.ClientConn, error) {
	return grpc.Dial(
		c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func (c Connection) NewClient() (*Client, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	return &Client{
		Cli:    pb.NewDatawayClient(conn),
		conn:   conn,
		logger: c.logger,
	}, nil
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		c.logger.Error(err)
	}
}
