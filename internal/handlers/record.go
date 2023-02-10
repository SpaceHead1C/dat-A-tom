package handlers

import (
	"context"
	"datatom/internal/api"
	"datatom/internal/domain"
	"errors"
	"fmt"
	"net/http"
)

func AddRecord(ctx context.Context, man *api.RecordManager, req AddRecordRequestSchema) (TextResult, error) {
	out := TextResult{Status: http.StatusCreated}
	r, err := req.AddRecordRequest()
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, err
	}
	id, err := man.Add(ctx, r)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = id.String()
	return out, nil
}

func UpdateRecord(ctx context.Context, man *api.RecordManager, req UpdRecordRequestSchema) (Result, error) {
	out := Result{Status: http.StatusNoContent}
	if req.Name == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("name %w", domain.ErrExpected)
	}
	if req.Description == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("description %w", domain.ErrExpected)
	}
	if req.DeletionMark == nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("deletion mark %w", domain.ErrExpected)
	}
	r, err := req.UpdRecordRequest()
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
