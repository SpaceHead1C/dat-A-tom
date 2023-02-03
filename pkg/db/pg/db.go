package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SSLMode uint

const (
	SSLDisable SSLMode = iota
	SSLRequire
	SSLVerifyCA
	SSLVerifyFull
)

func (m SSLMode) String() string {
	switch m {
	case SSLRequire:
		return "require"
	case SSLVerifyCA:
		return "verify-ca"
	case SSLVerifyFull:
		return "verify-full"
	}
	return "disable"
}

type Config struct {
	Address      string
	Port         uint
	User         string
	Password     string
	DatabaseName string
	SSLMode      SSLMode
}

func NewDB(ctx context.Context, c Config) (*pgx.Conn, error) {
	dbUrl := connectionString(c)
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectionString(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Address, c.Port, c.DatabaseName)
}
