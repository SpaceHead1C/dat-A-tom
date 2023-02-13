package handlers

import (
	"datatom/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

type AddPropertyRequestSchema struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Types          []string `json:"types"`
	RefTypeIDs     []string `json:"reference_type_ids"`
	OwnerRefTypeID string   `json:"owner_reference_type_id"`
}

func (s AddPropertyRequestSchema) AddPropertyRequest() (domain.AddPropertyRequest, []string, error) {
	out := domain.AddPropertyRequest{
		Name:        s.Name,
		Description: s.Description,
	}

	var unknownTypes []string
	tps := make([]domain.Type, 0, len(s.Types))
	var tp domain.Type
	for _, code := range s.Types {
		tp = domain.TypeFromCode(code)
		if tp == domain.UndefinedType {
			unknownTypes = append(unknownTypes, code)
		} else {
			tps = append(tps, tp)
		}
	}
	out.Types = tps

	rtps := make([]uuid.UUID, 0, len(s.RefTypeIDs))
	for _, v := range s.RefTypeIDs {
		id, err := uuid.Parse(v)
		if err != nil {
			return out, nil, fmt.Errorf("parse reference type id error: %s", err)
		}
		rtps = append(rtps, id)
	}
	out.RefTypeIDs = rtps

	if s.OwnerRefTypeID != "" {
	ortID, err := uuid.Parse(s.OwnerRefTypeID)
	if err != nil {
		return out, nil, fmt.Errorf("parse owner reference type id error: %s", err)
	}
	out.OwnerRefTypeID = ortID
	}

	if len(unknownTypes) > 0 {
		return out, unknownTypes, nil
	}
	return out, nil, nil
}

type UpdPropertyRequestSchema struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *UpdPropertyRequestSchema) UpdPropertyRequest() (domain.UpdPropertyRequest, error) {
	out := domain.UpdPropertyRequest{
		Name:        s.Name,
		Description: s.Description,
	}
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return out, fmt.Errorf("parse property id error: %s", err)
	}
	out.ID = id
	return out, nil
}
