package pg

import (
	"context"
	. "datatom/internal/domain"
	. "datatom/pkg/db/pg"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddProperty(ctx context.Context, req AddPropertyRequest) (uuid.UUID, error) {
	var out uuid.UUID
	args := []any{
		req.Name,
		req.Description,
		TypesToCodes(req.Types),
		ArrayUUID(req.RefTypes),
		NullUUID(req.OwnerRefType),
	}
	query := `SELECT new_property($1, $2, $3, $4, $5);`
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
