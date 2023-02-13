package handlers

import (
	"context"
	"datatom/internal/api"
	"datatom/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func UpdateProperty(ctx context.Context, man *api.PropertyManager, req UpdPropertyRequestSchema) (Result, error) {
	out := Result{Status: http.StatusNoContent}
	if req.Name == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("name %w", domain.ErrExpected)
	}
	if req.Description == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("description %w", domain.ErrExpected)
	}
	r, err := req.UpdPropertyRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	if _, err := man.Update(ctx, r); err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusNotFound
		}
		return out, err
	}
	return out, nil
}

func PatchProperty(ctx context.Context, man *api.PropertyManager, req UpdPropertyRequestSchema) (Result, error) {
	out := Result{Status: http.StatusOK}
	r, err := req.UpdPropertyRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	property, err := man.Update(ctx, r)
	if err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusNotFound
		}
		return out, err
	}
	b, err := json.Marshal(PropertyToResponseSchema(*property))
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}

func GetProperty(ctx context.Context, man *api.PropertyManager, id string) (Result, error) {
	out := Result{Status: http.StatusOK}
	rid, err := uuid.Parse(id)
	if err != nil {
		out.Status = http.StatusBadRequest
		out.Payload = []byte(fmt.Sprintf("parse property id error: %s", err))
		return out, err
	}
	property, err := man.Get(ctx, rid)
	if err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusNotFound
		}
		return out, err
	}
	b, err := json.Marshal(PropertyToResponseSchema(*property))
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}
