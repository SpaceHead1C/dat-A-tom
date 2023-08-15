package pb

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"net/http"
)

func UUIDToPb(in uuid.UUID) *UUID {
	return &UUID{Value: in[:]}
}

func UUIDFromPb(in *UUID) (uuid.UUID, error) {
	return uuid.FromBytes(in.Value)
}

func GRPCCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unavailable:
		return http.StatusMethodNotAllowed
	case codes.AlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
