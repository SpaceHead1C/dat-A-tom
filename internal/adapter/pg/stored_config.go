package pg

import (
	"context"
	. "datatom/internal/domain"
	"datatom/pkg/db/pg"
	"fmt"
	"github.com/google/uuid"
)

func (r *Repository) SetStoredConfigDatawayTomID(ctx context.Context, value uuid.UUID) error {
	query := `SELECT set_config_dataway_tom_id($1);`
	if _, err := r.ExecEx(ctx, query, nil, pg.NullUUID(value)); err != nil {
		return fmt.Errorf("database error: %w, %s\"", err, query)
	}
	return nil
}

func (r *Repository) GetStoredConfigDatawayTomID(ctx context.Context) (StoredConfigValue, error) {
	var id uuid.NullUUID
	out := StoredConfigUUID{Value: uuid.Nil}
	query := `SELECT get_config_dataway_tom_id();`
	if err := r.QueryRowEx(ctx, query, nil).Scan(&id); err != nil {
		return out, fmt.Errorf("database error: %w, %s\"", err, query)
	}
	if !id.Valid {
		return out, ErrStoredConfigTomIDNotSet
	}
	out.Value = id.UUID
	return out, nil
}
