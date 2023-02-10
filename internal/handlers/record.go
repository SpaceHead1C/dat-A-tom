package handlers

import (
	"context"
	"datatom/internal/api"
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
