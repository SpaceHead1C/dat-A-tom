package domain

import (
	"context"
	"fmt"
)

type ChangedDataType uint

type ChangedDataRepository interface {
	GetChanges(context.Context) ([]ChangedData, error)
}

type ChangedData struct {
	ID       int64
	DataType ChangedDataType
	Key      []byte
}

const (
	UnknownChangedDataType ChangedDataType = iota
	ChangedDataRefType
	ChangedDataProperty
	ChangedDataRecord
	ChangedDataValue
)

func (cdt ChangedDataType) Code() (string, error) {
	switch cdt {
	case ChangedDataRefType:
		return "ref_type", nil
	case ChangedDataProperty:
		return "property", nil
	case ChangedDataRecord:
		return "record", nil
	case ChangedDataValue:
		return "value", nil
	default:
		return "", fmt.Errorf("%w of changed data", ErrUnknownType)
	}
}

func (cdt ChangedDataType) String() string {
	switch cdt {
	case ChangedDataRefType:
		return "ref_type"
	case ChangedDataProperty:
		return "property"
	case ChangedDataRecord:
		return "record"
	case ChangedDataValue:
		return "value"
	default:
		return "unknown"
	}
}

func ChangedDataTypeFromCode(code string) ChangedDataType {
	switch code {
	case "ref_type":
		return ChangedDataRefType
	case "property":
		return ChangedDataProperty
	case "record":
		return ChangedDataRecord
	case "value":
		return ChangedDataValue
	default:
		return UnknownChangedDataType
	}
}
