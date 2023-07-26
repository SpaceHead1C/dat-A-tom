package api

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const defaultPropertyManagerTimeout = time.Second

type PropertyManager struct {
	PropertyConfig
}

type PropertyConfig struct {
	Repository PropertyRepository
	Broker     PropertyBroker
	Timeout    time.Duration
}

func NewPropertyManager(c PropertyConfig) (*PropertyManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("property repository can not be nil")
	}
	if c.Broker == nil {
		return nil, fmt.Errorf("property broker can not be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultPropertyManagerTimeout
	}
	return &PropertyManager{c}, nil
}

func (pm *PropertyManager) Add(ctx context.Context, req AddPropertyRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Repository.AddProperty(ctx, req)
}

func (pm *PropertyManager) Update(ctx context.Context, req UpdPropertyRequest) (*Property, error) {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Repository.UpdateProperty(ctx, req)
}

func (pm *PropertyManager) Get(ctx context.Context, id uuid.UUID) (*Property, error) {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Repository.GetProperty(ctx, id)
}

func (pm *PropertyManager) GetSentState(ctx context.Context, id uuid.UUID, transaction db.Transaction) (*PropertySentState, error) {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Repository.GetPropertySentStateForUpdate(ctx, id, transaction)
}

func (pm *PropertyManager) SetSentState(ctx context.Context, state PropertySentState, transaction db.Transaction) (*PropertySentState, error) {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Repository.SetSentProperty(ctx, state, transaction)
}

func (pm *PropertyManager) Send(ctx context.Context, req SendPropertyRequest) error {
	ctx, cancel := context.WithTimeout(ctx, pm.Timeout)
	defer cancel()
	return pm.Broker.SendProperty(ctx, req)
}

func (pm *PropertyManager) GetSender(req SendPropertyRequest) *PropertySender {
	return &PropertySender{
		man: pm,
		req: req,
	}
}
