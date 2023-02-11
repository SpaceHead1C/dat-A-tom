package helper

import "github.com/google/uuid"

func IsZeroUUID(id uuid.UUID) bool {
	return id == uuid.Nil
}
