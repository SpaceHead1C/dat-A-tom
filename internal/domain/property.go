package domain

import (
	"context"

	"github.com/google/uuid"
)

type PropertyRepository interface {
	AddProperty(context.Context, AddPropertyRequest) (uuid.UUID, error)
}

type AddPropertyRequest struct {
	Types          []Type
	RefTypeIDs     []uuid.UUID
	Name           string
	Description    string
	OwnerRefTypeID uuid.UUID
}
