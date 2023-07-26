package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type StoredConfig uint

const (
	StoredConfigTomID StoredConfig = iota
)

type storedConfigGetFunc func(StoredConfigRepository, context.Context) (StoredConfigValue, error)

func (sc StoredConfig) String() string {
	switch sc {
	case StoredConfigTomID:
		return "dat(A)way tom ID"
	default:
		return "unknown"
	}
}

func (sc StoredConfig) GetFunc() (storedConfigGetFunc, error) {
	switch sc {
	case StoredConfigTomID:
		return StoredConfigRepository.GetStoredConfigDatawayTomID, nil
	default:
		return nil, fmt.Errorf(`unexpected stored config "%s"`, sc.String())
	}
}

type StoredConfigRepository interface {
	SetStoredConfigDatawayTomID(context.Context, uuid.UUID) error
	GetStoredConfigDatawayTomID(context.Context) (StoredConfigValue, error)
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
		err = fmt.Errorf("%w %T", ErrUnexpectedType, dest)
	}
	return err
}
