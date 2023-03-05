package handlers

import "github.com/google/uuid"

type StoredConfigTomIDSchema struct {
	ID *string `json:"id"`
}

func TomIDToResponseSchema(id uuid.UUID, valid bool) StoredConfigTomIDSchema {
	var out StoredConfigTomIDSchema
	if valid {
		v := id.String()
		out.ID = &v
	}
	return out
}
