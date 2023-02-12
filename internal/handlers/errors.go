package handlers

import (
	. "datatom/internal/domain"
	"errors"
)

var badRequestErrors map[error]struct{}

func init() {
	badRequestErrors = map[error]struct{}{
		ErrTypesExpectedPG:            {},
		ErrTypesConditionNotMatchedPG: {},
		ErrTypeDuplicatedPG:           {},
		ErrRefTypeDuplicatedPG:        {},
		ErrUnknownRefTypePG:           {},
		ErrUnexpectedTypePG:           {},
		ErrUnexpectedRefTypePG:        {},
		ErrRefTypeExpectedPG:          {},
		ErrRefTypeIsRedundantPG:       {},
	}
}

func isBadRequestError(err error) bool {
	_, ok := badRequestErrors[err]
	return ok || errors.Is(err, ErrParseError)
}
