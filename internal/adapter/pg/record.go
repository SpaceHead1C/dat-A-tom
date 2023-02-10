package pg

import (
	"context"
	. "datatom/internal/domain"
	. "datatom/pkg/db/pg"
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
		NullUUID(req.ReferenceTypeID),
	}
	query := `SELECT new_record($1, $2, $3, $4);`
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

func (r *Repository) UpdateRecord(ctx context.Context, req UpdRecordRequest) (*Record, error) {
	args := make([]any, 4)
	args[0] = req.ID
	if req.Name != nil {
		args[1] = *req.Name
	}
	if req.Description != nil {
		args[2] = *req.Description
	}
	if req.DeletionMark != nil {
		args[3] = *req.DeletionMark
	}
	var recordJSON []byte
	query := `SELECT * FROM update_record($1, $2, $3, $4);`
	if err := r.QueryRowEx(ctx, query, nil, args...).Scan(&recordJSON); err != nil {
		if IsNoRowsError(err) {
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
