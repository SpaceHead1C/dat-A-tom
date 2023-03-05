package handlers

import (
	"context"
	"datatom/grpc"
	"datatom/internal/api"
	"datatom/internal/domain"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func RegisterTom(ctx context.Context, gRPCConn *grpc.Connection, man *api.StoredConfigsManager) (TextResult, error) {
	out := TextResult{Status: http.StatusCreated}
	if gRPCConn == nil {
		out.Status = http.StatusServiceUnavailable
		out.Payload = "disconnected from dat(A)way service"
		return out, nil
	}
	client, err := gRPCConn.NewClient()
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	defer client.Close()
	pbID, err := client.Cli.RegisterNewTom(ctx, &emptypb.Empty{})
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	id, err := uuid.FromBytes(pbID.Value)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	if err := man.Set(ctx, domain.StoredConfigTomID, id); err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = id.String()
	return out, nil
}

func GetTomID(ctx context.Context, man *api.StoredConfigsManager) (Result, error) {
	out := Result{Status: http.StatusOK}
	valid := true
	val, err := man.Get(ctx, domain.StoredConfigTomID)
	if err != nil {
		if errors.Is(err, domain.ErrStoredConfigTomIDNotSet) {
			err = nil
			valid = false
		} else {
			out.Status = http.StatusInternalServerError
			return out, err
		}
	}
	var id uuid.UUID
	if err := val.ScanStoredConfigValue(&id); err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	b, err := json.Marshal(TomIDToResponseSchema(id, valid))
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	out.Payload = b
	return out, nil
}
