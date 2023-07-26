package pg

import (
	"datatom/internal/domain"
	"github.com/google/uuid"
	"time"
)

type RefTypeSchema struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Sum         string    `json:"sum"`
	ChangeAt    time.Time `json:"change_at"`
}

func (rts *RefTypeSchema) RefType() *domain.RefType {
	return &domain.RefType{
		ID:          rts.ID,
		Name:        rts.Name,
		Description: rts.Description,
		Sum:         rts.Sum,
		ChangeAt:    rts.ChangeAt.UTC(),
	}
}
