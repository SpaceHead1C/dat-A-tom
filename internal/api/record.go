package api

import (
	"context"
	. "datatom/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const defaultRecordManagerTimeout = time.Second

type RecordManager struct {
	RecordConfig
}

type RecordConfig struct {
	Repository RecordRepository
	Timeout    time.Duration
}

func NewRecordManager(c RecordConfig) (*RecordManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("record repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultRecordManagerTimeout
	}
	return &RecordManager{c}, nil
}

func (rm *RecordManager) Add(ctx context.Context, req AddRecordRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Repository.AddRecord(ctx, req)
}

func (rm *RecordManager) Update(ctx context.Context, req UpdRecordRequest) (*Record, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Repository.UpdateRecord(ctx, req)
}

func (rm *RecordManager) Get(ctx context.Context, id uuid.UUID) (*Record, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Repository.GetRecord(ctx, id)
}
