package api

import (
	"context"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"fmt"
	"time"
)

const defaultChangedDataManagerTimeout = time.Second * 10

type ChangedDataManager struct {
	ChangedDataConfig
}

type ChangedDataConfig struct {
	Repository domain.ChangedDataRepository
	Timeout    time.Duration
}

func NewChangedDataManager(c ChangedDataConfig) (*ChangedDataManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("changed data repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultChangedDataManagerTimeout
	}
	return &ChangedDataManager{c}, nil
}

func (cdm *ChangedDataManager) Get(ctx context.Context) ([]domain.ChangedData, error) {
	ctx, cancel := context.WithTimeout(ctx, cdm.Timeout)
	defer cancel()
	return cdm.Repository.GetChanges(ctx)
}

func (cdm *ChangedDataManager) Purge(ctx context.Context, id int64, transaction db.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, cdm.Timeout)
	defer cancel()
	return cdm.Repository.PurgeChanges(ctx, id, transaction)
}
