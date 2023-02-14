package pg

import (
	. "datatom/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ValueSchema struct {
	OwnerID    uuid.UUID       `json:"owner_id"`
	PropertyID uuid.UUID       `json:"property_id"`
	Type       string          `json:"type"`
	RefTypeID  uuid.UUID       `json:"reference_type_id"`
	Val        ValueJSONSchema `json:"value"`
	Sum        string          `json:"sum"`
	ChangeAt   time.Time       `json:"change_at"`
}

func (vs *ValueSchema) Value() (*Value, error) {
	tp := TypeFromCode(vs.Type)
	value, err := ValidatedValue(vs.Val.V, tp)
	if err != nil {
		return nil, err
	}
	return &Value{
		RecordID:   vs.OwnerID,
		PropertyID: vs.PropertyID,
		Type:       tp,
		RefTypeID:  vs.RefTypeID,
		Value:      value,
		Sum:        vs.Sum,
		ChangeAt:   vs.ChangeAt.UTC(),
	}, nil
}
