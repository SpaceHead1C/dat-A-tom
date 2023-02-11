package pg

import (
	. "datatom/internal/domain"
	"time"

	"github.com/google/uuid"
)

type RecordSchema struct {
	ID              uuid.UUID `json:"id"`
	ReferenceTypeID uuid.UUID `json:"reference_type_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DeletionMark    bool      `json:"deletion_mark"`
	Sum             string    `json:"sum"`
	ChangeAt        time.Time `json:"change_at"`
}

func (rs *RecordSchema) Record() *Record {
	return &Record{
		ID:              rs.ID,
		ReferenceTypeID: rs.ReferenceTypeID,
		Name:            rs.Name,
		Description:     rs.Description,
		DeletionMark:    rs.DeletionMark,
		Sum:             rs.Sum,
		ChangeAt:        rs.ChangeAt.UTC(),
	}
}
