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
	Broker     ValueBroker
	Timeout    time.Duration
}

func NewValueManager(c ValueConfig) (*ValueManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("value repository can not be nil")
	}
	if c.Broker == nil {
		return nil, fmt.Errorf("value broker can not be nil")
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

func (vm *ValueManager) GetByKey(ctx context.Context, key []byte) (*Value, error) {
	req, err := getValueRequestByKey(key)
	if err != nil {
		return nil, fmt.Errorf("invalid key error: %w, %s", err, key)
	}
	return vm.Get(ctx, *req)
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

func (vm *ValueManager) Send(ctx context.Context, req SendValueRequest) error {
	ctx, cancel := context.WithTimeout(ctx, vm.Timeout)
	defer cancel()
	return vm.Broker.SendValue(ctx, req)
}

func (vm *ValueManager) GetSender(req SendValueRequest) *ValueSender {
	return &ValueSender{
		man: vm,
		req: req,
	}
}
