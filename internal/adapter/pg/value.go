package pg

import (
	"context"
	. "datatom/internal/domain"
	. "datatom/pkg/db/pg"
	"encoding/json"
	"fmt"
)

func (r *Repository) SetValue(ctx context.Context, req SetValueRequest) (*Value, error) {
	value, err := ValueAsJSON(req.Value, req.Type)
	if err != nil {
		return nil, err
	}
	var valueJSON []byte
	args := []any{
		NullUUID(req.RecordID),
		NullUUID(req.PropertyID),
		req.Type.Code(),
		NullUUID(req.RefTypeID),
		string(value),
	}
	query := `SELECT set_value($1, $2, $3, $4, $5);`
	if err := r.QueryRowEx(ctx, query, nil, args...).Scan(&valueJSON); err != nil {
		if errException, ok := pgExceptionAsDomainError(err); ok {
			return nil, errException
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema ValueSchema
	if err := json.Unmarshal(valueJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, valueJSON)
	}
	return schema.Value()
}
