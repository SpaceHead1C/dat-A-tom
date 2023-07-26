package pg

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	PoolMaxConn  uint
}

func NewConnConfig(c Config) (*pgx.ConnConfig, error) {
	return pgx.ParseConfig(connectionString(c))
}

func NewPoolConfig(c Config) (*pgxpool.Config, error) {
	if c.PoolMaxConn == 0 {
		c.PoolMaxConn = 10
	}
	return pgxpool.ParseConfig(poolConnectionString(c))
}

func connectionString(c Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=%s",
		c.User, c.Password, c.Address, c.Port, c.DatabaseName, c.SSLMode.String(),
	)
}

func poolConnectionString(c Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d",
		c.User, c.Password, c.Address, c.Port, c.DatabaseName, c.SSLMode.String(), c.PoolMaxConn,
	)
}
