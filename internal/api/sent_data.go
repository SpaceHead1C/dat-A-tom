package api

import (
	"context"
	. "datatom/internal/domain"
	"fmt"
	"time"
)

const defaultSentDataManagerTimeout = time.Second

type SentDataManager struct {
	SentDataConfig
}

type SentDataConfig struct {
	Repository SentDataRepository
	Timeout    time.Duration
}

func NewSentDataManager(c SentDataConfig) (*SentDataManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("sent data repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultSentDataManagerTimeout
	}
	return &SentDataManager{c}, nil
}

func (sdm *SentDataManager) Get(ctx context.Context, req GetSentDataRequest) (*SentData, error) {
	ctx, cancel := context.WithTimeout(ctx, sdm.Timeout)
	defer cancel()
	return sdm.Repository.GetSentData(ctx, req)
}
