package rest

import (
	"datatom/internal/handlers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newAddPropertyHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.AddPropertyRequestSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		res, err := handlers.AddProperty(req.Context(), s.propertyManager, schema)
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("add property error: %s", err)
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

func newUpdPropertyHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.UpdPropertyRequestSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		schema.ID = chi.URLParam(req, "id")
		res, err := handlers.UpdateProperty(req.Context(), s.propertyManager, schema)
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("update property error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.emptyResp(w, res.Status)
	})
}

func newPatchPropertyHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.UpdPropertyRequestSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		schema.ID = chi.URLParam(req, "id")
		res, err := handlers.PatchProperty(req.Context(), s.propertyManager, schema)
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("patch property error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.jsonResp(w, res.Status, res.Payload)
	})
}

func newGetPropertyHandler(s *server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.GetProperty(req.Context(), s.propertyManager, chi.URLParam(req, "id"))
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("get property error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.jsonResp(w, res.Status, res.Payload)
	})
}
