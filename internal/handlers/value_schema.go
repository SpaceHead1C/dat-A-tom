package handlers

import (
	"datatom/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

type SetValueRequestSchema struct {
	RecordID   string `json:"record_id"`
	PropertyID string `json:"property_id"`
	Type       string `json:"type"`
	RefTypeID  string `json:"reference_type_id"`
	Value      any    `json:"value"`
}

func (s SetValueRequestSchema) SetValueRequest() (domain.SetValueRequest, error) {
	out := domain.SetValueRequest{
		Value: s.Value,
	}
	recordID, err := uuid.Parse(s.RecordID)
	if err != nil {
		return out, fmt.Errorf("parse record id error: %s", err)
	}
	propertyID, err := uuid.Parse(s.PropertyID)
	if err != nil {
		return out, fmt.Errorf("parse property id error: %s", err)
	}
	tp := domain.TypeFromCode(s.Type)
	if tp == domain.UndefinedType {
		return out, fmt.Errorf("unknown type %s", s.Type)
	}
	if s.RefTypeID != "" {
		id, err := uuid.Parse(s.RefTypeID)
		if err != nil {
			return out, fmt.Errorf("parse reference type id error: %s", err)
		}
		out.RefTypeID = id
	}
	out.RecordID = recordID
	out.PropertyID = propertyID
	out.Type = tp
	return out, nil
}
