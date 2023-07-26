package pg

import (
	"context"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"datatom/pkg/db/pg"
	"fmt"
)

func (r *Repository) GetChanges(ctx context.Context) ([]domain.ChangedData, error) {
	var out []domain.ChangedData
	var id int64
	var dataType string
	var key []byte
	query := `SELECT * FROM get_changes();`
	rows, err := r.Query(ctx, query)
	if err != nil {
		if pg.IsNoRowsError(err) {
			return out, nil
		}
		return nil, fmt.Errorf("database error: %w, %s", err, query)
	}
	for rows.Next() {
		if err := rows.Scan(&id, &dataType, &key); err != nil {
			return nil, fmt.Errorf("database scan error: %w, %s", err, query)
		}
		out = append(out, domain.ChangedData{
			ID:       id,
			DataType: domain.ChangedDataTypeFromCode(dataType),
			Key:      key,
		})
	}
	return out, nil
}

func (r *Repository) PurgeChanges(ctx context.Context, id int64, tx db.Transaction) error {
	query := `SELECT purge_changes($1);`
	queryRow, err := funcQueryRow(r, tx)
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}
	var deleted int64
	if err := queryRow(ctx, query, id).Scan(&deleted); err != nil {
		return fmt.Errorf("database error: %w, %s", err, query)
	}
	return nil
}
