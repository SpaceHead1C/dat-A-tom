package handlers

import (
	"context"
	"datatom/grpc"
	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/internal/pb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type SubscribeRequest struct {
	GRPCConn    *grpc.Connection
	PropertyMan *api.PropertyManager
	SCMan       *api.StoredConfigsManager
	Schema      SubscribeSchema
}

type DeleteSubscriptionRequest struct {
	GRPCConn    *grpc.Connection
	PropertyMan *api.PropertyManager
	SCMan       *api.StoredConfigsManager
	Schema      SubscribeSchema
}

func RegisterTom(ctx context.Context, gRPCConn *grpc.Connection, man *api.StoredConfigsManager) (TextResult, error) {
	out := TextResult{Status: http.StatusCreated}
	if gRPCConn == nil {
		out.Status = http.StatusMethodNotAllowed
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
	id, err := pb.UUIDFromPb(pbID)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, fmt.Errorf("read ID of registered tom error: %w", err)
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
	id, valid, err := getTomID(ctx, man)
	if err != nil {
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

func Subscribe(ctx context.Context, req SubscribeRequest) (Result, error) {
	out := Result{Status: http.StatusCreated}
	tomID, valid, err := getTomID(ctx, req.SCMan)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	if !valid {
		out.Status = http.StatusMethodNotAllowed
		return out, fmt.Errorf("tom not registered in dat(A)way service")
	}
	propertyID, err := uuid.Parse(req.Schema.PropertyID)
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("parse property ID error: %s", err)
	}
	if _, err := req.PropertyMan.Get(ctx, propertyID); err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusBadRequest
		}
		return out, err
	}
	consumerID, err := uuid.Parse(req.Schema.ConsumerID)
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("parse consumer ID error: %s", err)
	}
	client, err := req.GRPCConn.NewClient()
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	defer client.Close()
	if _, err := client.Cli.Subscribe(ctx, &pb.Subscription{
		ConsumerId: pb.UUIDToPb(consumerID),
		TomId:      pb.UUIDToPb(tomID),
		PropertyId: pb.UUIDToPb(propertyID),
	}); err != nil {
		out.Status = http.StatusInternalServerError
		if s, ok := status.FromError(err); ok {
			out.Status = pb.GRPCCodeToHTTPStatus(s.Code())
			err = fmt.Errorf(s.Message())
		}
		return out, err
	}
	return out, nil
}

func DeleteSubscription(ctx context.Context, req DeleteSubscriptionRequest) (Result, error) {
	out := Result{Status: http.StatusOK}
	tomID, valid, err := getTomID(ctx, req.SCMan)
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	if !valid {
		out.Status = http.StatusMethodNotAllowed
		return out, fmt.Errorf("tom not registered in dat(A)way service")
	}
	propertyID, err := uuid.Parse(req.Schema.PropertyID)
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("parse property ID error: %s", err)
	}
	if _, err := req.PropertyMan.Get(ctx, propertyID); err != nil {
		out.Status = http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			out.Status = http.StatusBadRequest
		}
		return out, err
	}
	consumerID, err := uuid.Parse(req.Schema.ConsumerID)
	if err != nil {
		out.Status = http.StatusBadRequest
		return out, fmt.Errorf("parse consumer ID error: %s", err)
	}
	client, err := req.GRPCConn.NewClient()
	if err != nil {
		out.Status = http.StatusInternalServerError
		return out, err
	}
	defer client.Close()
	if _, err := client.Cli.DeleteSubscription(ctx, &pb.Subscription{
		ConsumerId: pb.UUIDToPb(consumerID),
		TomId:      pb.UUIDToPb(tomID),
		PropertyId: pb.UUIDToPb(propertyID),
	}); err != nil {
		out.Status = http.StatusInternalServerError
		if s, ok := status.FromError(err); ok {
			out.Status = pb.GRPCCodeToHTTPStatus(s.Code())
			err = fmt.Errorf(s.Message())
		}
		return out, err
	}
	return out, nil
}

func getTomID(ctx context.Context, man *api.StoredConfigsManager) (id uuid.UUID, valid bool, err error) {
	valid = true
	val, err := man.Get(ctx, domain.StoredConfigTomID)
	if err != nil {
		if errors.Is(err, domain.ErrStoredConfigTomIDNotSet) {
			err = nil
			valid = false
		} else {
			return
		}
	}
	err = val.ScanStoredConfigValue(&id)
	return
}
