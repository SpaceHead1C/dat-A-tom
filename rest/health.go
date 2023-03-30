package rest

import (
	"datatom/internal/handlers"
	"net/http"
)

func newPingHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res := handlers.Ping()
		s.emptyResp(w, res.Status)
	}
}

func newInfoHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.Info(s.appInfo)
		if err != nil {
			s.logger.Errorf("get app info error: %s", err)
			s.emptyResp(w, res.Status)
		}
		s.jsonResp(w, res.Status, res.Payload)
	}
}
