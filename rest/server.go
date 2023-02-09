package rest

import (
	"datatom/internal/domain"
	"datatom/pkg/log"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const (
	defaultHTTPServerTimeout = time.Second * 5
)

type server struct {
	logger         *zap.SugaredLogger
	srv            *http.Server
	errorHandler   func(error)
	timeout        time.Duration
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

type Config struct {
	Logger         *zap.SugaredLogger
	Port           uint
	ErrorHandler   func(error)
	Timeout        time.Duration
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
	if c.Timeout == 0 {
		c.Timeout = defaultHTTPServerTimeout
	}
	out := &server{
		logger:         l,
		errorHandler:   eh,
		timeout:        c.Timeout,
	}

	router := chi.NewRouter()
	router.Use(mw.StripSlashes)
	router.Use(mw.GetHead)
	router.Use(mw.Timeout(out.timeout))

	router.Mount("/health", healthRouter(out))

	out.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		WriteTimeout: time.Second * 7,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}
	return out, nil
}

func healthRouter(s *server) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", newPingHandler(s))
	return r
}
