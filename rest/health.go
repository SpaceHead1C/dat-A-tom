package rest

import (
	"datatom/internal/handlers"
	"net/http"
)

func newPingHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res := handlers.Ping()
		s.emptyResp(w, res.Status)
	})
}
