package rest

import (
	"datatom/internal/api"
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

	regexUUIDTemplate = `[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`
)

type server struct {
	logger         *zap.SugaredLogger
	srv            *http.Server
	errorHandler   func(error)
	timeout        time.Duration
	refTypeManager *api.RefTypeManager
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

type Config struct {
	Logger         *zap.SugaredLogger
	Port           uint
	ErrorHandler   func(error)
	Timeout        time.Duration
	RefTypeManager *api.RefTypeManager
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
	if c.RefTypeManager == nil {
		return nil, fmt.Errorf("reference type manager must not be nil")
	}
	if c.Timeout == 0 {
		c.Timeout = defaultHTTPServerTimeout
	}
	out := &server{
		logger:         l,
		errorHandler:   eh,
		timeout:        c.Timeout,
		refTypeManager: c.RefTypeManager,
	}

	router := chi.NewRouter()
	router.Use(mw.StripSlashes)
	router.Use(mw.GetHead)
	router.Use(mw.Timeout(out.timeout))

	router.Mount("/health", healthRouter(out))
	router.Mount("/ref_type", refTypeRouter(out))

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

func refTypeRouter(s *server) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", newAddRefTypeHandler(s))
	r.Put(fmt.Sprintf("/{id:%s}", regexUUIDTemplate), newUpdRefTypeHandler(s))
	r.Patch(fmt.Sprintf("/{id:%s}", regexUUIDTemplate), newPatchRefTypeHandler(s))
	r.Get(fmt.Sprintf("/{id:%s}", regexUUIDTemplate), newGetRefTypeHandler(s))
	return r
}

func (s *server) emptyResp(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func (s *server) textResp(w http.ResponseWriter, status int, payload string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if err := writeResp(w, status, []byte(payload)); err != nil {
		s.errorHandler(err)
	}
}

func (s *server) jsonResp(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := writeResp(w, status, payload); err != nil {
		s.errorHandler(err)
	}
}

func writeResp(w http.ResponseWriter, status int, payload []byte) error {
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return nil
}
