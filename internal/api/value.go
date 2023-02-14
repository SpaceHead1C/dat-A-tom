package api

import (
	"context"
	. "datatom/internal/domain"
	"fmt"
	"time"
)

const defaultValueManagerTimeout = time.Second

type ValueManager struct {
	ValueConfig
}

type ValueConfig struct {
	Repository ValueRepository
	Timeout    time.Duration
}

func NewValueManager(c ValueConfig) (*ValueManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("value repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultValueManagerTimeout
	}
	return &ValueManager{c}, nil
}

func (vm *ValueManager) Set(ctx context.Context, req SetValueRequest) (*Value, error) {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Repository.SetValue(ctx, req)
}
