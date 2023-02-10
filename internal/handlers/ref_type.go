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

func UpdateRefType(ctx context.Context, man *api.RefTypeManager, req UpdRefTypeRequestSchema) (Result, error) {
	out := Result{Status: http.StatusNoContent}
	if req.Name == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("name %w", domain.ErrExpected)
	}
	if req.Description == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("description %w", domain.ErrExpected)
	}
	r, err := req.UpdRefTypeRequest()
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

func PatchRefType(ctx context.Context, man *api.RefTypeManager, req UpdRefTypeRequestSchema) (Result, error) {
	out := Result{Status: http.StatusOK}
	r, err := req.UpdRefTypeRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	refType, err := man.Update(ctx, r)
	if err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusNotFound
		}
		return out, err
	}
	b, err := json.Marshal(RefTypeToResponseSchema(*refType))
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}

func GetRefType(ctx context.Context, man *api.RefTypeManager, id string) (Result, error) {
	out := Result{Status: http.StatusOK}
	rid, err := uuid.Parse(id)
	if err != nil {
		out.Status = http.StatusBadRequest
		out.Payload = []byte(fmt.Sprintf("parse reference type id error: %s", err))
		return out, err
	}
	refType, err := man.Get(ctx, rid)
	if err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusNotFound
		}
		return out, err
	}
	b, err := json.Marshal(RefTypeToResponseSchema(*refType))
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}
