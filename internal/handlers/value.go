package handlers

import (
	"context"
	"datatom/internal/api"
	"net/http"
)

func SetValue(ctx context.Context, man *api.ValueManager, req SetValueRequestSchema) (Result, error) {
	out := Result{Status: http.StatusNoContent}
	r, err := req.SetValueRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	if _, err := man.Set(ctx, r); err != nil {
		out.Status = http.StatusInternalServerError
		if isBadRequestError(err) {
			out.Status = http.StatusBadRequest
		}
		return out, err
	}
	return out, nil
}
