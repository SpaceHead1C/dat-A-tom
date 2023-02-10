package handlers

import (
	"datatom/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

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

type UpdRefTypeRequestSchema struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *UpdRefTypeRequestSchema) UpdRefTypeRequest() (domain.UpdRefTypeRequest, error) {
	out := domain.UpdRefTypeRequest{
		Name:        s.Name,
		Description: s.Description,
	}
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return out, fmt.Errorf("parse reference type id error: %s", err)
	}
	out.ID = id
	return out, nil
}

type RefTypeResponseSchema struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func RefTypeToResponseSchema(rt domain.RefType) RefTypeResponseSchema {
	return RefTypeResponseSchema{
		ID:          rt.ID.String(),
		Name:        rt.Name,
		Description: rt.Description,
	}
}
