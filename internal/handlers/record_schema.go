package handlers

import (
	"datatom/internal/domain"
	"datatom/pkg/helper"
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
		Name:         s.Name,
		Description:  s.Description,
		DeletionMark: s.DeletionMark,
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

type UpdRecordRequestSchema struct {
	ID           string  `json:"id"`
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DeletionMark *bool   `json:"deletion_mark,omitempty"`
}

func (s UpdRecordRequestSchema) UpdRecordRequest() (domain.UpdRecordRequest, error) {
	out := domain.UpdRecordRequest{
		Name:         s.Name,
		Description:  s.Description,
		DeletionMark: s.DeletionMark,
	}
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return out, fmt.Errorf("parse record id error: %s", err)
	}
	out.ID = id
	return out, nil
}

type RecordResponseSchema struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DeletionMark    bool    `json:"deletion_mark"`
	ReferenceTypeID *string `json:"reference_type_id"`
}

func RecordToResponseSchema(rt domain.Record) RecordResponseSchema {
	var refTypeID *string
	if !helper.IsZeroUUID(rt.ReferenceTypeID) {
		rtID := rt.ReferenceTypeID.String()
		refTypeID = &rtID
	}
	return RecordResponseSchema{
		ID:              rt.ID.String(),
		Name:            rt.Name,
		Description:     rt.Description,
		DeletionMark:    rt.DeletionMark,
		ReferenceTypeID: refTypeID,
	}
}
