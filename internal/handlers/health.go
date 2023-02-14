package handlers

import (
	"net/http"
)

func Ping() Result {
	return Result{Status: http.StatusOK}
}
