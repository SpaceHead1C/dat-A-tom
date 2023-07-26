package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db"
	"datatom/pkg/db/pg"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

func (r *Repository) SetSentRecord(ctx context.Context, state RecordSentState, tx db.Transaction) (*RecordSentState, error) {
	var sentDataJSON []byte
	args := []any{
		pg.NullUUID(state.ID),
		state.Sum,
		state.SentAt,
	}
	query := `SELECT set_sent_record($1, $2, $3);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RecordSentStateSchema
	if err := json.Unmarshal(sentDataJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentDataJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) GetRecordSentStateForUpdate(ctx context.Context, id uuid.UUID, tx db.Transaction) (*RecordSentState, error) {
	var sentStateJSON []byte
	query := `SELECT get_sent_record_for_update($1);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, id).Scan(&sentStateJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrSentDataNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RecordSentStateSchema
	if err := json.Unmarshal(sentStateJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentStateJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) SetSentProperty(ctx context.Context, state PropertySentState, tx db.Transaction) (*PropertySentState, error) {
	var sentDataJSON []byte
	args := []any{
		pg.NullUUID(state.ID),
		state.Sum,
		state.SentAt,
	}
	query := `SELECT set_sent_property($1, $2, $3);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema PropertySentStateSchema
	if err := json.Unmarshal(sentDataJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentDataJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) GetPropertySentStateForUpdate(ctx context.Context, id uuid.UUID, tx db.Transaction) (*PropertySentState, error) {
	var sentStateJSON []byte
	query := `SELECT get_sent_property_for_update($1);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, id).Scan(&sentStateJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrSentDataNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema PropertySentStateSchema
	if err := json.Unmarshal(sentStateJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentStateJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) SetSentRefType(ctx context.Context, state RefTypeSentState, tx db.Transaction) (*RefTypeSentState, error) {
	var sentDataJSON []byte
	args := []any{
		pg.NullUUID(state.ID),
		state.Sum,
		state.SentAt,
	}
	query := `SELECT set_sent_ref_type($1, $2, $3);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RefTypeSentStateSchema
	if err := json.Unmarshal(sentDataJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentDataJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) GetRefTypeSentStateForUpdate(ctx context.Context, id uuid.UUID, tx db.Transaction) (*RefTypeSentState, error) {
	var sentStateJSON []byte
	query := `SELECT get_sent_ref_type_for_update($1);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, id).Scan(&sentStateJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrSentDataNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema RefTypeSentStateSchema
	if err := json.Unmarshal(sentStateJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentStateJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) SetSentValue(ctx context.Context, state ValueSentState, tx db.Transaction) (*ValueSentState, error) {
	var sentDataJSON []byte
	args := []any{
		pg.NullUUID(state.RecordID),
		pg.NullUUID(state.PropertyID),
		state.Sum,
		state.SentAt,
	}
	query := `SELECT set_sent_value($1, $2, $3, $4);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema ValueSentStateSchema
	if err := json.Unmarshal(sentDataJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentDataJSON)
	}
	return schema.SentData(), nil
}

func (r *Repository) GetValueSentStateForUpdate(ctx context.Context, req GetValueRequest, tx db.Transaction) (*ValueSentState, error) {
	var sentStateJSON []byte
	args := []any{
		req.RecordID,
		req.PropertyID,
	}
	query := `SELECT get_sent_value_for_update($1, $2);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}
	if err := queryRow(ctx, query, args...).Scan(&sentStateJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrSentDataNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema ValueSentStateSchema
	if err := json.Unmarshal(sentStateJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentStateJSON)
	}
	return schema.SentData(), nil
}
