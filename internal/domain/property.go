package domain

import (
	"context"

	"github.com/google/uuid"
)

type PropertyRepository interface {
	AddProperty(context.Context, AddPropertyRequest) (uuid.UUID, error)
}

type AddPropertyRequest struct {
	Types        []Type
	RefTypes     []uuid.UUID
	Code         string
	Name         string
	Description  string
	OwnerRefType uuid.UUID
}
