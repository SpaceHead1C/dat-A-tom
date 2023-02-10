package domain

import (
	"context"

	"github.com/google/uuid"
)

type RecordRepository interface {
	AddRecord(context.Context, AddRecordRequest) (uuid.UUID, error)
}

type AddRecordRequest struct {
	Name          string
	Description   string
	DeletionMark  bool
	ReferenceType uuid.UUID
}
