package api

import (
	"context"
	. "datatom/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (rtm *RefTypeManager) Add(req AddRefTypeRequest) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rtm.Timeout)
	defer cancel()
	return rtm.Repository.AddRefType(ctx, req)
}