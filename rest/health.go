package rest

import (
	"datatom/internal/handlers"
	"io"
	"net/http"
)

func newPingHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, handlers.Ping())
		if err != nil {
			s.errorHandler(err)
		}
	})
}
