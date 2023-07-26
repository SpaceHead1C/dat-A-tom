package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrRecordNotFound   = fmt.Errorf("record %w", ErrNotFound)
	ErrRefTypeNotFound  = fmt.Errorf("reference type %w", ErrNotFound)
	ErrPropertyNotFound = fmt.Errorf("property %w", ErrNotFound)
	ErrValueNotFound    = fmt.Errorf("value %w", ErrNotFound)
	ErrSentDataNotFound = fmt.Errorf("sent data %w", ErrNotFound)

	ErrStoredConfigTomIDNotSet = errors.New("tom ID not set")

	ErrExpected       = errors.New("expected")
	ErrUnexpectedType = errors.New("unexpected type")
	ErrUnknownType    = errors.New("unknown type")
	ErrParseError     = errors.New("parse")

	// PostgreSQL exceptions
	ErrTypesExpectedPG            = errors.New("types expected")
	ErrTypesConditionNotMatchedPG = errors.New("types and reference type condition not matched")
	ErrTypeDuplicatedPG           = errors.New("type duplicated")
	ErrRefTypeDuplicatedPG        = errors.New("reference type duplicated")
	ErrUnknownRefTypePG           = errors.New("unknown reference type")
	ErrUnexpectedTypePG           = fmt.Errorf("%w", ErrUnexpectedType)
	ErrUnexpectedRefTypePG        = errors.New("unexpected reference type")
	ErrRefTypeExpectedPG          = errors.New("reference type expected")
	ErrRefTypeIsRedundantPG       = errors.New("no need reference type ID cause type is not reference")
)
