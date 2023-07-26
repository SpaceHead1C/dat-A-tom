package api

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db"
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

func (vm *ValueManager) Get(ctx context.Context, req GetValueRequest) (*Value, error) {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Repository.GetValue(ctx, req)
}

func (vm *ValueManager) GetSentState(ctx context.Context, req GetValueRequest, transaction db.Transaction) (*ValueSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Repository.GetValueSentStateForUpdate(ctx, req, transaction)
}

func (vm *ValueManager) SetSentState(ctx context.Context, state ValueSentState, transaction db.Transaction) (*ValueSentState, error) {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Repository.SetSentValue(ctx, state, transaction)
}

func (vm *ValueManager) ChangedValues(ctx context.Context) ([]Value, error) {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Repository.ChangedValues(ctx)
}
