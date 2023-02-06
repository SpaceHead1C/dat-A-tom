package pg

import (
	"datatom/pkg/db/pg"
	"fmt"

	"github.com/jackc/pgx/stdlib"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

const getUUIDAttemptsThreshold = 10

type Repository struct {
	*pgx.Conn
	l *zap.SugaredLogger
}

func NewRepository(db *pg.DB, l *zap.SugaredLogger) (*Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	conn, err := stdlib.AcquireConn(db.DB)
	if err != nil {
		return nil, err
	}
	return &Repository{conn, l}, nil
}

func (r *Repository) CloseConn(db *pg.DB) {
	stdlib.ReleaseConn(db.DB, r.Conn)
}
