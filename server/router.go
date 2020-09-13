package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gsiragusa/short-to-me/errors"
	"github.com/sirupsen/logrus"
)

// RouteHandler adds error handling to the normal http.Handler interface
type RouteHandler func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r)
func (f RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// RouteConfig describes the configuration of a handler for a HTTP route
type RouteConfig struct {
	Name    string
	Method  string
	Path    string
	Handler RouteHandler
}

type Router struct {
	MuxRouter *mux.Router
}

// NewRouter returns a mux.Router for the api server
func NewRouter(routeConfigArr ...[]RouteConfig) *Router {
	r := &Router{
		MuxRouter: mux.NewRouter().StrictSlash(true),
	}

	// all handlers
	for _, routeConfigs := range routeConfigArr {
		for _, routeConfig := range routeConfigs {
			r.addHandler(routeConfig)
		}
	}

	// not found handler
	r.MuxRouter.NotFoundHandler = r.notFoundHandler()

	// public documentation
	r.MuxRouter.PathPrefix("/docs/").Handler(
		http.StripPrefix("/docs/",
			http.FileServer(http.Dir("public_docs/"))))

	return r
}

func (r *Router) addHandler(routes ...RouteConfig) {
	for _, routeConfig := range routes {
		route := r.MuxRouter.
			Methods(routeConfig.Method).
			Path(routeConfig.Path).
			Name(routeConfig.Name)

		route.Handler(httpHandler(routeConfig.Handler))
	}
}

func (r *Router) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = WriteError(w, errors.NewErrorNotFound())
	})
}

// httpHandler wraps all handlers
func httpHandler(h RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := h.ServeHTTP(w, req)
		errHandler(w, err)
	}
}

// errHandler ensure that errors are handled properly
func errHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	_ = WriteError(w, errors.NewInternalServerError())
}

func Write(w http.ResponseWriter, code int, payload interface{}) error {
	// short circuit if the payload is empty
	if payload == nil {
		w.WriteHeader(code)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// must WriteHeader after the potential error, writing before would create
	// a panic when failing to marshal the payload, as router catches the error
	// and attempts to write a new status to the response
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		// no point passing this error on as something has gone wrong in the
		// http layer at this point and beyond the applications control to send
		// an error response - a retry will be performed by http
		logrus.WithError(err).Error("failed to write response")
	}

	return nil
}

func WriteError(w http.ResponseWriter, err errors.Error) error {
	return Write(w, err.HttpStatus, errors.ErrResponse(err))
}
