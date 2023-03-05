package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type StoredConfig uint

const (
	StoredConfigTomID StoredConfig = iota
)

func (sc StoredConfig) String() string {
	switch sc {
	case StoredConfigTomID:
		return "dat(A)way tom ID"
	default:
		return "unknown"
	}
}

type StoredConfigValue interface {
	ScanStoredConfigValue(dest any) error
}

type StoredConfigUUID struct {
	Value uuid.UUID
}

func (s StoredConfigUUID) ScanStoredConfigValue(dest any) error {
	if dest == nil {
		return nil
	}
	var err error
	if d, ok := dest.(*uuid.UUID); ok {
		*d = s.Value
		dest = d
	} else {
		err = fmt.Errorf("unexpected type %T", dest)
	}
	return err
}
