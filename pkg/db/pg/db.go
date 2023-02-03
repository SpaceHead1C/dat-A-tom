package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
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

type DB struct {
	*sql.DB
}

func NewDB(ctx context.Context, c Config) (*DB, error) {
	db, err := sql.Open("pgx", connectionString(c))
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func connectionString(c Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=%s",
		c.User, c.Password, c.Address, c.Port, c.DatabaseName, c.SSLMode.String(),
	)
}
