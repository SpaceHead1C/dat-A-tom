package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddRecord(ctx context.Context, req AddRecordRequest) (uuid.UUID, error) {
	var out uuid.UUID
	args := []any{
		req.Name,
		req.Description,
		req.DeletionMark,
		pg.NullUUID(req.ReferenceTypeID),
	}
	query := `SELECT new_record($1, $2, $3, $4);`
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

func (r *Repository) UpdateRecord(ctx context.Context, req UpdRecordRequest) (*Record, error) {
	emptyReq := true
	args := make([]any, 4)
	args[0] = req.ID
	if req.Name != nil {
		args[1] = *req.Name
		emptyReq = false
	}
	if req.Description != nil {
		args[2] = *req.Description
		emptyReq = false
	}
	if req.DeletionMark != nil {
		args[3] = *req.DeletionMark
		emptyReq = false
	}
	if emptyReq {
		return r.GetRecord(ctx, req.ID)
	}
	var recordJSON []byte
	query := `SELECT * FROM update_record($1, $2, $3, $4);`
	if err := r.QueryRow(ctx, query, args...).Scan(&recordJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RecordSchema
	if err := json.Unmarshal(recordJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, recordJSON)
	}
	return schema.Record(), nil
}

func (r *Repository) GetRecord(ctx context.Context, id uuid.UUID) (*Record, error) {
	var recordJSON []byte
	query := `SELECT * FROM get_record($1);`
	if err := r.QueryRow(ctx, query, id).Scan(&recordJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RecordSchema
	if err := json.Unmarshal(recordJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, recordJSON)
	}
	return schema.Record(), nil
}
