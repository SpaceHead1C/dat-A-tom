package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) SetSentData(ctx context.Context, req SetSentDataRequest) (*SentData, error) {
	var sentDataJSON []byte
	args := []any{
		req.RecordID,
		req.PropertyID,
		req.Sum,
		req.SentAt,
	}
	r.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	})
	query := `SELECT set_sent_data($1, $2, $3, $4);`
	if err := r.QueryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema SentDataSchema
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
		req.RecordID,
		req.PropertyID,
	}
	query := `SELECT get_sent_data($1, $2);`
	if err := r.QueryRow(ctx, query, args...).Scan(&sentDataJSON); err != nil {
		if pg.IsNoRowsError(err) {
			return nil, ErrSentDataNotFound
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	var schema SentDataSchema
	if err := json.Unmarshal(sentDataJSON, &schema); err != nil {
		return nil, fmt.Errorf("db result unmarshal error: %s, %s", err, sentDataJSON)
	}
	return schema.SentData(), nil
}
