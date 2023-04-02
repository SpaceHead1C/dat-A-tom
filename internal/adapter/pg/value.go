package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
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
		pg.NullUUID(req.RecordID),
		pg.NullUUID(req.PropertyID),
		req.Type.Code(),
		pg.NullUUID(req.RefTypeID),
		string(value),
	}
	query := `SELECT set_value($1, $2, $3, $4, $5);`
	if err := r.QueryRow(ctx, query, args...).Scan(&valueJSON); err != nil {
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

func (r *Repository) ChangedValues(ctx context.Context) ([]Value, error) {
	query := `SELECT * FROM get_changed_values();`
	rows, err := r.Query(ctx, query)
	if err != nil {
		if pg.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var out []Value
	for rows.Next() {
		var valueJSON []byte
		if err := rows.Scan(&valueJSON); err != nil {
			return nil, fmt.Errorf("database scan error: %w, %s", err, query)
		}
		var schema ValueSchema
		if err := json.Unmarshal(valueJSON, &schema); err != nil {
			return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, valueJSON)
		}
		value, err := schema.Value()
		if err != nil {
			return nil, err
		}
		out = append(out, *value)
	}
	return out, nil
}
