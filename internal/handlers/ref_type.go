package handlers

import (
	"context"
	"datatom/internal/api"
	"net/http"
)

func AddRefType(ctx context.Context, man *api.RefTypeManager, req AddRefTypeRequestSchema) (TextResult, error) {
	out := TextResult{Status: http.StatusCreated}
	id, err := man.Add(ctx, req.AddRefTypeRequest())
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = id.String()
	return out, nil
}
