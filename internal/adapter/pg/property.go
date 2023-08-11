package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
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
		pg.ArrayUUID(req.RefTypeIDs),
		pg.NullUUID(req.OwnerRefTypeID),
	}
	query := `SELECT new_property($1, $2, $3, $4, $5);`
	for attempts := 0; attempts < getUUIDAttemptsThreshold; attempts++ {
		if err := r.QueryRow(ctx, query, args...).Scan(&out); err != nil {
			if pg.IsNotUniqueError(err) {
				continue
			}
			if errException, ok := pgExceptionAsDomainError(err); ok {
				return out, errException
			}
			return out, fmt.Errorf("database error: %w, %s", err, query)
		}
		return out, nil
	}
	return out, errCanNotGetUniqueID
}

func (r *Repository) UpdateProperty(ctx context.Context, req UpdPropertyRequest) (*Property, error) {
	emptyReq := true
	args := make([]any, 3)
	args[0] = req.ID
	if req.Name != nil {
		args[1] = *req.Name
		emptyReq = false
	}
	if req.Description != nil {
		args[2] = *req.Description
		emptyReq = false
	}
	if emptyReq {
		return r.GetProperty(ctx, req.ID)
	}
	var propertyJSON []byte
	query := `SELECT update_property($1, $2, $3);`
	if err := r.QueryRow(ctx, query, args...).Scan(&propertyJSON); err != nil {
		if pg.IsNoRowsError(err) {
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
	if err := r.QueryRow(ctx, query, id).Scan(&propertyJSON); err != nil {
		if pg.IsNoRowsError(err) {
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
