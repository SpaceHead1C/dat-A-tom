package handlers

import "datatom/internal/domain"

type AddRefTypeRequestSchema struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s AddRefTypeRequestSchema) AddRefTypeRequest() domain.AddRefTypeRequest {
	return domain.AddRefTypeRequest{
		Name:        s.Name,
		Description: s.Description,
	}
}
