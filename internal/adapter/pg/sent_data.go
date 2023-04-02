package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
	"encoding/json"
	"fmt"
)

func (r *Repository) GetSentData(ctx context.Context, req GetSentDataRequest) (*SentData, error) {
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
