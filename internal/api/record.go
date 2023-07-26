package api

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db"
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
	Broker     RecordBroker
	Timeout    time.Duration
}

func NewRecordManager(c RecordConfig) (*RecordManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("record repository can not be nil")
	}
	if c.Broker == nil {
		return nil, fmt.Errorf("record broker can not be nil")
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

func (rm *RecordManager) GetSentState(ctx context.Context, id uuid.UUID, transaction db.Transaction) (*RecordSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Repository.GetRecordSentStateForUpdate(ctx, id, transaction)
}

func (rm *RecordManager) SetSentState(ctx context.Context, state RecordSentState, transaction db.Transaction) (*RecordSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Repository.SetSentRecord(ctx, state, transaction)
}

func (rm *RecordManager) Send(ctx context.Context, req SendRecordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, rm.Timeout)
	defer cancel()
	return rm.Broker.SendRecord(ctx, req)
}

func (rm *RecordManager) GetSender(req SendRecordRequest) *RecordSender {
	return &RecordSender{
		man: rm,
		req: req,
	}
}
