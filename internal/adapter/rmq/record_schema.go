package rmq

import (
	"datatom/internal/domain"
	"datatom/pkg/helper"
)

type RecordSchema struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DeletionMark    bool    `json:"deletion_mark"`
	ReferenceTypeID *string `json:"reference_type_id"`
}

func recordToSchema(r domain.Record) RecordSchema {
	var refTypeID *string
	if !helper.IsZeroUUID(r.ReferenceTypeID) {
		rtID := r.ReferenceTypeID.String()
		refTypeID = &rtID
	}
	return RecordSchema{
		ID:              r.ID.String(),
		Name:            r.Name,
		Description:     r.Description,
		DeletionMark:    r.DeletionMark,
		ReferenceTypeID: refTypeID,
	}
}
