package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddRefType(ctx context.Context, req AddRefTypeRequest) (uuid.UUID, error) {
	var out uuid.UUID
	args := []any{
		req.Name,
		req.Description,
	}
	query := `SELECT new_ref_type($1, $2);`
	for attempts := 0; attempts < getUUIDAttemptsThreshold; attempts++ {
		if err := r.QueryRow(ctx, query, args...).Scan(&out); err != nil {
			if pg.IsNotUniqueError(err) {
				continue
			}
			return out, fmt.Errorf("database error: %w, %s", err, query)
		}
		return out, nil
	}
	return out, errCanNotGetUniqueID
}

func (r *Repository) UpdateRefType(ctx context.Context, req UpdRefTypeRequest) (*RefType, error) {
	var out RefType
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
		return r.GetRefType(ctx, req.ID)
	}
	query := `SELECT * FROM update_ref_type($1, $2, $3);`
	if err := r.QueryRow(ctx, query, args...).Scan(&out.ID, &out.Name, &out.Description); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrRefTypeNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	return &out, nil
}

func (r *Repository) GetRefType(ctx context.Context, id uuid.UUID) (*RefType, error) {
	var refTypeJSON []byte
	query := `SELECT get_ref_type($1);`
	if err := r.QueryRow(ctx, query, id).Scan(&refTypeJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrRefTypeNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RefTypeSchema
	if err := json.Unmarshal(refTypeJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, refTypeJSON)
	}
	return schema.RefType(), nil
}
