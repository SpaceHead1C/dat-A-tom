package rest

import (
	"datatom/internal/domain"
	"datatom/pkg/log"
	"fmt"
	"net/http"
	"time"
	"go.uber.org/zap"
)

type server struct {
	logger       *zap.SugaredLogger
	srv          *http.Server
	errorHandler func(error)
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

type Config struct {
	Logger       *zap.SugaredLogger
	Port         int
	ErrorHandler func(error)
}

func NewServer(c Config) (domain.Server, error) {
	var err error
	l := c.Logger
	if l == nil {
		l, err = log.NewLogger()
		if err != nil {
			return nil, err
		}
	}
	eh := c.ErrorHandler
	if eh == nil {
		eh = func(e error) {
			l.Errorln(e.Error())
		}
	}
	out := &server{
		logger:       l,
		errorHandler: eh,
	}

	out.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		WriteTimeout: time.Second * 2,
		ReadTimeout:  time.Second * 9,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}
	return out, nil
}
