package pg

import (
	. "datatom/internal/domain"
	"fmt"

	"github.com/jackc/pgx"
)

var (
	pgExceptions map[string]error

	errCanNotGetUniqueID = fmt.Errorf("can not get unique ID")
)

func init() {
	pgExceptions = map[string]error{
		"types expected": ErrTypesExpectedPG,
		"types and reference type condition not matched": ErrTypesConditionNotMatchedPG,
		"type duplicated":                                       ErrTypeDuplicatedPG,
		"reference type ID duplicated":                          ErrRefTypeDuplicatedPG,
		"unknown reference type ID":                             ErrUnknownRefTypePG,
		"unexpected type":                                       ErrUnexpectedTypePG,
		"unexpected reference type ID":                          ErrUnexpectedRefTypePG,
		"reference type ID missing":                             ErrRefTypeExpectedPG,
		"no need reference type ID cause type is not reference": ErrRefTypeIsRedundantPG,
	}
}

func pgExceptionAsDomainError(err error) (error, bool) {
	if err == nil {
		return nil, false
	}
	pgErr, ok := err.(pgx.PgError)
	if !ok {
		return err, false
	}
	if pgErr.Code != "P0001" {
		return err, false
	}
	if errException, ok := pgExceptions[pgErr.Message]; ok {
		return errException, true
	}
	return err, false
}
