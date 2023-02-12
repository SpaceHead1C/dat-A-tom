package api

import (
	"context"
	. "datatom/internal/domain"
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
	Timeout    time.Duration
}

func NewPropertyManager(c PropertyConfig) (*PropertyManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("property repository can't be nil")
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
