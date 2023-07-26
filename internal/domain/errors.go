package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = fmt.Errorf("not found")
	ErrRecordNotFound   = fmt.Errorf("record %w", ErrNotFound)
	ErrRefTypeNotFound  = fmt.Errorf("reference type %w", ErrNotFound)
	ErrPropertyNotFound = fmt.Errorf("property %w", ErrNotFound)
	ErrSentDataNotFound = fmt.Errorf("sent data %w", ErrNotFound)

	ErrUnexpectedType          = errors.New("unexpected type")
	ErrStoredConfigTomIDNotSet = errors.New("tom ID not set")

	ErrExpected    = errors.New("expected")
	ErrUnknownType = errors.New("unknown type")
	ErrParseError  = errors.New("parse")

	// PostgreSQL exceptions
	ErrTypesExpectedPG            = fmt.Errorf("types expected")
	ErrTypesConditionNotMatchedPG = fmt.Errorf("types and reference type condition not matched")
	ErrTypeDuplicatedPG           = fmt.Errorf("type duplicated")
	ErrRefTypeDuplicatedPG        = fmt.Errorf("reference type duplicated")
	ErrUnknownRefTypePG           = fmt.Errorf("unknown reference type")
	ErrUnexpectedTypePG           = fmt.Errorf("unexpected type")
	ErrUnexpectedRefTypePG        = fmt.Errorf("unexpected reference type")
	ErrRefTypeExpectedPG          = fmt.Errorf("reference type expected")
	ErrRefTypeIsRedundantPG       = fmt.Errorf("no need reference type ID cause type is not reference")
)
