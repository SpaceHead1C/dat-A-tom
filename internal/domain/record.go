package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RecordRepository interface {
	AddRecord(context.Context, AddRecordRequest) (uuid.UUID, error)
	UpdateRecord(context.Context, UpdRecordRequest) (*Record, error)
}

type Record struct {
	ID              uuid.UUID
	ReferenceTypeID uuid.UUID
	Name            string
	Description     string
	DeletionMark    bool
	Sum             string
	ChangeAt        time.Time
}

type AddRecordRequest struct {
	Name            string
	Description     string
	DeletionMark    bool
	ReferenceTypeID uuid.UUID
}

type UpdRecordRequest struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	DeletionMark *bool
}
