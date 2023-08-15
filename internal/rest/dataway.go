package rest

import (
	"datatom/internal/handlers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newRegisterTomHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.RegisterTom(req.Context(), handlers.RegisterTomRequest{
			GRPCConn: s.dwGRPCConn,
			SCMan:    s.storedConfigsManager,
			AppInfo:  s.appInfo,
		})
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest, http.StatusConflict, http.StatusMethodNotAllowed:
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

func newUpdateTomHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.UpdateTom(req.Context(), handlers.RegisterTomRequest{
			GRPCConn: s.dwGRPCConn,
			SCMan:    s.storedConfigsManager,
			AppInfo:  s.appInfo,
		})
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest, http.StatusConflict, http.StatusMethodNotAllowed:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("register tom error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.emptyResp(w, res.Status)
	}
}

func newGetTomIDHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := handlers.GetTomID(req.Context(), s.storedConfigsManager)
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
	}
}

func newSubscribeHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.SubscribeSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		res, err := handlers.Subscribe(req.Context(), handlers.SubscribeRequest{
			GRPCConn:    s.dwGRPCConn,
			PropertyMan: s.propertyManager,
			SCMan:       s.storedConfigsManager,
			Schema:      schema,
		})
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest, http.StatusMethodNotAllowed:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("add subscription error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.emptyResp(w, res.Status)
	}
}

func newDeleteSubscriptionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			s.textResp(w, http.StatusInternalServerError, "read body error")
			s.logger.Errorf("read body error: %s", err)
			return
		}
		var schema handlers.SubscribeSchema
		if err := json.Unmarshal(b, &schema); err != nil {
			s.textResp(w, http.StatusBadRequest, fmt.Sprintf("body unmarshal error: %s", err))
			return
		}
		res, err := handlers.DeleteSubscription(req.Context(), handlers.DeleteSubscriptionRequest{
			GRPCConn:    s.dwGRPCConn,
			PropertyMan: s.propertyManager,
			SCMan:       s.storedConfigsManager,
			Schema:      schema,
		})
		if err != nil {
			switch res.Status {
			case http.StatusBadRequest, http.StatusMethodNotAllowed:
				s.textResp(w, res.Status, err.Error())
			case http.StatusInternalServerError:
				s.logger.Errorf("delete subscription error: %s", err)
				fallthrough
			default:
				s.emptyResp(w, res.Status)
			}
			return
		}
		s.emptyResp(w, res.Status)
	}
}
