package pg

import (
	. "datatom/internal/domain"
	"time"

	"github.com/google/uuid"
)

type PropertySchema struct {
	ID             uuid.UUID   `json:"id"`
	Types          []string    `json:"types"`
	RefTypeIDs     []uuid.UUID `json:"reference_type_ids"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	OwnerRefTypeID uuid.UUID   `json:"owner_reference_type_id"`
	Sum            string      `json:"sum"`
	ChangeAt       time.Time   `json:"change_at"`
}

func (rs *PropertySchema) Property() *Property {
	return &Property{
		ID:             rs.ID,
		Name:           rs.Name,
		Description:    rs.Description,
		Types:          TypesFromCodes(rs.Types),
		RefTypeIDs:     rs.RefTypeIDs,
		OwnerRefTypeID: rs.OwnerRefTypeID,
		Sum:            rs.Sum,
		ChangeAt:       rs.ChangeAt.UTC(),
	}
}
