package api

import (
	"context"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const defaultChangedDataManagerTimeout = time.Second * 10

type ChangedDataManager struct {
	ChangedDataConfig
}

type ChangedDataConfig struct {
	Repository domain.ChangedDataRepository
	Timeout    time.Duration
}

func NewChangedDataManager(c ChangedDataConfig) (*ChangedDataManager, error) {
	if c.Repository == nil {
		return nil, fmt.Errorf("changed data repository can't be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultChangedDataManagerTimeout
	}
	return &ChangedDataManager{c}, nil
}

func (cdm *ChangedDataManager) Get(ctx context.Context) ([]domain.ChangedData, error) {
	ctx, cancel := context.WithTimeout(ctx, cdm.Timeout)
	defer cancel()
	return cdm.Repository.GetChanges(ctx)
}

func (cdm *ChangedDataManager) Purge(ctx context.Context, id int64, transaction db.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, cdm.Timeout)
	defer cancel()
	return cdm.Repository.PurgeChanges(ctx, id, transaction)
}

type changedValueKeySchema struct {
	OwnerID    *uuid.UUID `json:"owner_id,omitempty"`
	PropertyID *uuid.UUID `json:"property_id,omitempty"`
}

func getValueRequestByKey(key []byte) (*domain.GetValueRequest, error) {
	var schema changedValueKeySchema
	if err := json.Unmarshal(key, &schema); err != nil {
		return nil, err
	}
	if schema.OwnerID == nil {
		return nil, errors.New("owner id is not present in schema")
	}
	if schema.PropertyID == nil {
		return nil, errors.New("property id is not present in schema")
	}
	return &domain.GetValueRequest{
		RecordID:   *schema.OwnerID,
		PropertyID: *schema.PropertyID,
	}, nil
}

type changedDataKeySchema struct {
	ID *uuid.UUID `json:"id,omitempty"`
}

func getDataRequestByKey(key []byte) (uuid.UUID, error) {
	var schema changedDataKeySchema
	if err := json.Unmarshal(key, &schema); err != nil {
		return uuid.Nil, err
	}
	if schema.ID == nil {
		return uuid.Nil, errors.New("id is not present in schema")
	}
	return *schema.ID, nil
}
