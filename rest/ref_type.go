package rest

import (
	"datatom/internal/handlers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newAddRefTypeHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.AddRefTypeRequestSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		res, err := handlers.AddRefType(req.Context(), s.refTypeManager, schema)
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("add reference type error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s/%s", req.URL.String(), res.Payload))
		s.textResp(w, res.Status, res.Payload)
	})
}