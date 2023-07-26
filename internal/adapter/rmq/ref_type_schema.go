package rmq

import "datatom/internal/domain"

type RefTypeResponseSchema struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func refTypeToSchema(rt domain.RefType) RefTypeResponseSchema {
	return RefTypeResponseSchema{
		ID:          rt.ID.String(),
		Name:        rt.Name,
		Description: rt.Description,
	}
}
