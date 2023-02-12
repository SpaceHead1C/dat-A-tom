package handlers

import (
	"context"
	"datatom/internal/api"
	"datatom/internal/domain"
	"fmt"
	"net/http"
)

func AddProperty(ctx context.Context, man *api.PropertyManager, req AddPropertyRequestSchema) (TextResult, error) {
	out := TextResult{Status: http.StatusCreated}
	r, unknownTypes, err := req.AddPropertyRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	if len(unknownTypes) > 0 {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("unknown types: %v", unknownTypes)
	}
	if len(r.Types) == 0 {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("types %w", domain.ErrExpected)
	}
	id, err := man.Add(ctx, r)
	if err != nil {
		out.Status = http.StatusInternalServerError
		if isBadRequestError(err) {
			out.Status = http.StatusBadRequest
		}
		return out, err
	}
	out.Payload = id.String()
	return out, nil
}
