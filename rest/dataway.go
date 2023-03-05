package rest

import (
	"datatom/internal/handlers"
	"net/http"
)

func newRegisterTomHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.RegisterTom(req.Context(), s.dwGRPCConn, s.storedConfigsManager)
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("register tom error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.textResp(w, res.Status, res.Payload)
	}
}
