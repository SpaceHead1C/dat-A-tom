package domain

import (
	"context"
	"datatom/pkg/db"
	"time"

	"github.com/google/uuid"
)

const DeliveryTypeProperty = "property"

type PropertyRepository interface {
	AddProperty(context.Context, AddPropertyRequest) (uuid.UUID, error)
	UpdateProperty(context.Context, UpdPropertyRequest) (*Property, error)
	GetProperty(context.Context, uuid.UUID) (*Property, error)
	GetPropertySentStateForUpdate(context.Context, uuid.UUID, db.Transaction) (*PropertySentState, error)
	SetSentProperty(context.Context, PropertySentState, db.Transaction) (*PropertySentState, error)
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

type PropertySentState struct {
	ID     uuid.UUID
	Sum    string
	SentAt time.Time
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
