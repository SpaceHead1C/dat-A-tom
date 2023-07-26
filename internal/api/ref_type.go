package api

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const defaultRefTypeManagerTimeout = time.Second

type RefTypeManager struct {
	RefTypeConfig
}

type RefTypeConfig struct {
	Repository RefTypeRepository
	Broker     RefTypeBroker
	Timeout    time.Duration
}

func NewRefTypeManager(c RefTypeConfig) (*RefTypeManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("reference type repository can not be nil")
	}
	if c.Broker == nil {
		return nil, fmt.Errorf("reference type broker can not be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultRefTypeManagerTimeout
	}
	return &RefTypeManager{c}, nil
}

func (rtm *RefTypeManager) Add(ctx context.Context, req AddRefTypeRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.AddRefType(ctx, req)
}

func (rtm *RefTypeManager) Update(ctx context.Context, req UpdRefTypeRequest) (*RefType, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.UpdateRefType(ctx, req)
}

func (rtm *RefTypeManager) Get(ctx context.Context, id uuid.UUID) (*RefType, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.GetRefType(ctx, id)
}

func (rtm *RefTypeManager) GetByKey(ctx context.Context, key []byte) (*RefType, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.GetRefTypeByKey(ctx, key)
}

func (rtm *RefTypeManager) GetSentState(ctx context.Context, id uuid.UUID, transaction db.Transaction) (*RefTypeSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.GetRefTypeSentStateForUpdate(ctx, id, transaction)
}

func (rtm *RefTypeManager) SetSentState(ctx context.Context, state RefTypeSentState, transaction db.Transaction) (*RefTypeSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Repository.SetSentRefType(ctx, state, transaction)
}

func (rtm *RefTypeManager) Send(ctx context.Context, req SendRefTypeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, rtm.Timeout)
	defer cancel()
	return rtm.Broker.SendRefType(ctx, req)
}

func (rtm *RefTypeManager) GetSender(req SendRefTypeRequest) *RefTypeSender {
	return &RefTypeSender{
		man: rtm,
		req: req,
	}
}
