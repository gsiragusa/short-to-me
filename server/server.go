package server

import (
	"fmt"
	"net/http"

	"github.com/gsiragusa/short-to-me/config"
	"github.com/ory/graceful"
	"github.com/sirupsen/logrus"
)

type Server struct {
	le     *logrus.Logger
	config *config.AppConfig
	router *Router
}

type Handler interface {
	GetRoutes() []RouteConfig
}

func New(le *logrus.Logger, appConfig *config.AppConfig, routes ...Handler) *Server {
	parsedRoutes := parseHandlers(routes)
	return &Server{
		le:     le,
		config: appConfig,
		router: NewRouter(parsedRoutes...),
	}
}

func (s *Server) ListenAndServe() error {
	s.le.Infof("starting api on %d", s.config.Port)
	srv := s.newServer()
	return graceful.Graceful(srv.ListenAndServe, srv.Shutdown)
}

func (s *Server) newServer() *http.Server {
	return graceful.WithDefaults(&http.Server{
		Handler: s.router.MuxRouter,
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		// timeouts can be set here
	})
}

func parseHandlers(handlers []Handler) [][]RouteConfig {
	var routes [][]RouteConfig
	for _, h := range handlers {
		routes = append(routes, h.GetRoutes())
	}
	return routes
}
