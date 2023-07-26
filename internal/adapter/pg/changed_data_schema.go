package pg

import (
	"datatom/internal/domain"
	"encoding/json"
	"github.com/google/uuid"
)

type changedValueKeySchema struct {
	OwnerID    uuid.UUID `json:"owner_id"`
	PropertyID uuid.UUID `json:"property_id"`
}

func getValueRequestByKey(key []byte) (*domain.GetValueRequest, error) {
	var schema changedValueKeySchema
	if err := json.Unmarshal(key, &schema); err != nil {
		return nil, err
	}
	return &domain.GetValueRequest{
		RecordID:   schema.OwnerID,
		PropertyID: schema.PropertyID,
	}, nil
}

type changedDataKeySchema struct {
	ID uuid.UUID `json:"id"`
}

func getDataRequestByKey(key []byte) (uuid.UUID, error) {
	var schema changedDataKeySchema
	if err := json.Unmarshal(key, &schema); err != nil {
		return uuid.Nil, err
	}
	return schema.ID, nil
}
