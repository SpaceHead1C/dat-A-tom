package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const getUUIDAttemptsThreshold = 10

type Config struct {
	ConnectConfig *pgx.ConnConfig
	Logger        *zap.SugaredLogger
}

type Repository struct {
	*pgx.Conn
	l *zap.SugaredLogger
}

func NewRepository(ctx context.Context, c Config) (*Repository, error) {
	if c.ConnectConfig == nil {
		return nil, fmt.Errorf("connect config is nil")
	}
	if c.Logger == nil {
		c.Logger = zap.L().Sugar()
	}
	conn, err := pgx.ConnectConfig(ctx, c.ConnectConfig)
	if err != nil {
		return nil, err
	}
	return &Repository{conn, c.Logger}, nil
}
