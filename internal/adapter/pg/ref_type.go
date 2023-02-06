package pg

import (
	"context"
	. "datatom/internal/domain"
	. "datatom/pkg/db/pg"
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
	if err := r.QueryRowEx(ctx, query, nil, args...).Scan(&out); err != nil {
		if IsNotUniqueError(err) {
			return out, errCanNotGetUniqueID
		}
		return out, fmt.Errorf("database error: %w, %s", err, query)
	}
	return out, nil
}
