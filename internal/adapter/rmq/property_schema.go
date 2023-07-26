package rmq

import (
	"datatom/internal/domain"
	"datatom/pkg/helper"
)

type PropertySchema struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Types          []string `json:"types"`
	RefTypeIDs     []string `json:"reference_type_ids"`
	OwnerRefTypeID *string  `json:"owner_reference_type_id"`
}

func propertyToSchema(p domain.Property) PropertySchema {
	var refTypeIDs []string
	if len(p.RefTypeIDs) > 0 {
		refTypeIDs = make([]string, 0, len(p.RefTypeIDs))
		for _, rtID := range p.RefTypeIDs {
			refTypeIDs = append(refTypeIDs, rtID.String())
		}
	}
	var ownerRefTypeID *string
	if !helper.IsZeroUUID(p.OwnerRefTypeID) {
		ortID := p.OwnerRefTypeID.String()
		ownerRefTypeID = &ortID
	}
	return PropertySchema{
		ID:             p.ID.String(),
		Name:           p.Name,
		Description:    p.Description,
		Types:          domain.TypesToCodes(p.Types),
		RefTypeIDs:     refTypeIDs,
		OwnerRefTypeID: ownerRefTypeID,
	}
}
