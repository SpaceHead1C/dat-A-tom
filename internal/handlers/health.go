package handlers

import (
	"datatom/internal"
	"encoding/json"
	"net/http"
)

func Ping() Result {
	return Result{Status: http.StatusOK}
}

func Info(i internal.Info) (Result, error) {
	out := Result{Status: http.StatusOK}
	schema := InfoResponseSchema{
		Service:     internal.ServiceName,
		Version:     i.Version(),
		Title:       i.Title(),
		Description: i.Description(),
	}
	b, err := json.Marshal(schema)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}
