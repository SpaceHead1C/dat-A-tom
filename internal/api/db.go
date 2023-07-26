package api

import (
	"context"
	"datatom/pkg/db"
	"fmt"
	"time"
)

const defaultDBManagerTimeout = time.Second * 30

type DBManager struct {
	DBConfig
}

type DBConfig struct {
	Repository db.TransactionBeginner
	Timeout    time.Duration
}

func NewDBManager(c DBConfig) (*DBManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("DB repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultDBManagerTimeout
	}
	return &DBManager{c}, nil
}

func (dbm *DBManager) BeginTransaction(ctx context.Context) (db.Transaction, error) {
	return dbm.Repository.BeginTransaction(ctx)
}
