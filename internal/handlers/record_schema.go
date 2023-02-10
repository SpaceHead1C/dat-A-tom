package handlers

import (
	"datatom/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

type AddRecordRequestSchema struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DeletionMark    bool   `json:"deletion_mark"`
	ReferenceTypeID string `json:"reference_type_id"`
}

func (s AddRecordRequestSchema) AddRecordRequest() (domain.AddRecordRequest, error) {
	out := domain.AddRecordRequest{
		Name:            s.Name,
		Description:     s.Description,
		DeletionMark:    s.DeletionMark,
	}
	if s.ReferenceTypeID != "" {
		id, err := uuid.Parse(s.ReferenceTypeID)
		if err != nil {
			return out, fmt.Errorf("parse reference type id error: %s", err)
		}
		out.ReferenceTypeID = id
	}
	return out, nil
}
