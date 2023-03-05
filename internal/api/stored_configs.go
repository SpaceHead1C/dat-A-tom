package api

import (
	"context"
	. "datatom/internal/domain"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const defaultStoredConfigsManagerTimeout = time.Second

type StoredConfigsManager struct {
	StoredConfigsConfig
}

type StoredConfigsConfig struct {
	Repository StoredConfigRepository
	Timeout    time.Duration
}

func NewStoredConfigManager(c StoredConfigsConfig) (*StoredConfigsManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("stored configs repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultStoredConfigsManagerTimeout
	}
	return &StoredConfigsManager{c}, nil
}

func (sm StoredConfigsManager) Set(ctx context.Context, sc StoredConfig, value any) error {
	ctx, cancel := context.WithTimeout(ctx, sm.Timeout)
	defer cancel()
	switch sc {
	case StoredConfigTomID:
		if id, ok := value.(uuid.UUID); ok {
			return sm.Repository.SetStoredConfigDatawayTomID(ctx, id)
		} else {
			return fmt.Errorf("unexpected type %T", value)
		}
	default:
		return fmt.Errorf(`unexpected stored config "%s"`, sc.String())
	}
}

func (sm StoredConfigsManager) Get(ctx context.Context, sc StoredConfig) (StoredConfigValue, error) {
	get, err := sc.GetFunc()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, sm.Timeout)
	defer cancel()
	return get(sm.Repository, ctx)
}
