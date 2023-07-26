package pg

import (
	"encoding/json"
	"github.com/google/uuid"
)

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
