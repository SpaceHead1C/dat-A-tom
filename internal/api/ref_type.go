package api

import (
	"fmt"
	"time"
)

const defaultRefTypeManagerTimeout = time.Second * 10

type RefTypeManager struct {
	RefTypeConfig
}

type RefTypeConfig struct {
	Repository RefTypeRepository
	Timeout    time.Duration
}

func NewRefTypeManager(c RefTypeConfig) (*RefTypeManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("reference type repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultRefTypeManagerTimeout
	}
	return &RefTypeManager{c}, nil
}
