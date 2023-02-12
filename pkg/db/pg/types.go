package pg

import (
	"datatom/pkg/helper"

	"github.com/google/uuid"
)

func NullUUID(v uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: v, Valid: !helper.IsZeroUUID(v)}
}

func ArrayUUID(in []uuid.UUID) [][16]byte {
	if in == nil {
		return nil
	}
	out := make([][16]byte, 0, len(in))
	for _, v := range in {
		out = append(out, [16]byte(v))
	}
	return out
}
