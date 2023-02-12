package pg

import (
	"context"
	. "datatom/internal/domain"
	. "datatom/pkg/db/pg"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddProperty(ctx context.Context, req AddPropertyRequest) (uuid.UUID, error) {
	var out uuid.UUID
	args := []any{
		req.Name,
		req.Description,
		TypesToCodes(req.Types),
		ArrayUUID(req.RefTypeIDs),
		NullUUID(req.OwnerRefTypeID),
	}
	query := `SELECT new_property($1, $2, $3, $4, $5);`
	for attempts := 0; attempts < getUUIDAttemptsThreshold; attempts++ {
		if err := r.QueryRowEx(ctx, query, nil, args...).Scan(&out); err != nil {
			if IsNotUniqueError(err) {
				continue
			}
			return out, fmt.Errorf("database error: %w, %s", err, query)
		}
		return out, nil
	}
	return out, errCanNotGetUniqueID
}

func (r *Repository) UpdateProperty(ctx context.Context, req UpdPropertyRequest) (*Property, error) {
	args := make([]any, 3)
	args[0] = req.ID
	if req.Name != nil {
		args[1] = *req.Name
	}
	if req.Description != nil {
		args[2] = *req.Description
	}
	var propertyJSON []byte
	query := `SELECT update_property($1, $2, $3);`
	if err := r.QueryRowEx(ctx, query, nil, args...).Scan(&propertyJSON); err != nil {
		if IsNoRowsError(err) {
			return nil, ErrPropertyNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema PropertySchema
	if err := json.Unmarshal(propertyJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, propertyJSON)
	}
	return schema.Property(), nil
}

func (r *Repository) GetProperty(ctx context.Context, id uuid.UUID) (*Property, error) {
	var propertyJSON []byte
	query := `SELECT get_property($1);`
	if err := r.QueryRowEx(ctx, query, nil, id).Scan(&propertyJSON); err != nil {
		if IsNoRowsError(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema PropertySchema
	if err := json.Unmarshal(propertyJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, propertyJSON)
	}
	return schema.Property(), nil
}
