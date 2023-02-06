package pg

import "github.com/jackc/pgx"

func IsNotUniqueError(err error) bool {
	if err == nil {
		return false
	}
	pgErr, ok := err.(pgx.PgError)
	return ok && pgErr.Code == "23505"
}
