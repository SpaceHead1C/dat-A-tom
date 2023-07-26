package rmq

import (
	"datatom/internal/domain"
	"datatom/pkg/helper"
)

type ValueSchema struct {
	RecordID        string  `json:"record_id"`
	PropertyID      string  `json:"property_id"`
	Type            string  `json:"type"`
	ReferenceTypeID *string `json:"reference_type_id"`
	Value           any     `json:"value"`
}

func valueToSchema(v domain.Value) ValueSchema {
	var referenceTypeID *string
	if !helper.IsZeroUUID(v.RefTypeID) {
		rtID := v.RefTypeID.String()
		referenceTypeID = &rtID
	}
	return ValueSchema{
		RecordID:        v.RecordID.String(),
		PropertyID:      v.PropertyID.String(),
		Type:            v.Type.Code(),
		ReferenceTypeID: referenceTypeID,
		Value:           v.Value,
	}
}
