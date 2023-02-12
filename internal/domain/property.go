package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PropertyRepository interface {
	AddProperty(context.Context, AddPropertyRequest) (uuid.UUID, error)
	UpdateProperty(context.Context, UpdPropertyRequest) (*Property, error)
	GetProperty(context.Context, uuid.UUID) (*Property, error)
}

type Property struct {
	ID             uuid.UUID
	Types          []Type
	RefTypeIDs     []uuid.UUID
	Name           string
	Description    string
	OwnerRefTypeID uuid.UUID
	Sum            string
	ChangeAt       time.Time
}

type AddPropertyRequest struct {
	Types          []Type
	RefTypeIDs     []uuid.UUID
	Name           string
	Description    string
	OwnerRefTypeID uuid.UUID
}

type UpdPropertyRequest struct {
	ID          uuid.UUID
	Name        *string
	Description *string
}
