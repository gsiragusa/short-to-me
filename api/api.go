package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/errors"
	"github.com/gsiragusa/short-to-me/server"
	"github.com/gsiragusa/short-to-me/shortener"
	"github.com/sirupsen/logrus"
)

type API struct {
	le   *logrus.Logger
	conf *config.AppConfig
	svc  shortener.Service
}

func NewAPI(le *logrus.Logger, conf *config.AppConfig, svc shortener.Service) *API {
	return &API{
		le:   le,
		conf: conf,
		svc:  svc,
	}
}

// GetRoutes defines the router paths handled by this API
func (api *API) GetRoutes() []server.RouteConfig {
	return []server.RouteConfig{
		{
			Name:    "create-short-url",
			Method:  http.MethodPost,
			Path:    "/api/shortener",
			Handler: api.createShortUrl,
		},
		{
			Name:    "read-short-url",
			Method:  http.MethodGet,
			Path:    "/api/shortener",
			Handler: api.readShortUrl,
		},
		{
			Name:    "delete-short-url",
			Method:  http.MethodDelete,
			Path:    "/api/shortener",
			Handler: api.deleteShortUrl,
		},
		{
			Name:    "count-redirects",
			Method:  http.MethodGet,
			Path:    "/api/count",
			Handler: api.countRedirects,
		},
		{
			Name:    "redirect",
			Method:  http.MethodGet,
			Path:    "/{shortId}",
			Handler: api.redirect,
		},
	}
}

// TODO document
func (api *API) createShortUrl(w http.ResponseWriter, r *http.Request) error {
	// parse and validate input
	url, err := api.parseInput(r)
	if err != nil {
		return server.WriteError(w, *err)
	}

	encoded, err := api.svc.ShortenUrl(r.Context(), url)
	if err != nil {
		return server.WriteError(w, *err)
	}

	shortUrl := fmt.Sprintf("http://%s/%s", r.Host, encoded)
	return server.Write(w, http.StatusOK, &shortener.ModelShorten{Url: shortUrl})
}

func (api *API) readShortUrl(w http.ResponseWriter, r *http.Request) error {
	// parse and validate input
	url, err := api.parseInput(r)
	if err != nil {
		return server.WriteError(w, *err)
	}

	res, err := api.svc.RetrieveUrl(r.Context(), url)
	if err != nil {
		return server.WriteError(w, errors.NewErrorNotFound())
	}

	return server.Write(w, http.StatusOK, &shortener.ModelShorten{Url: res})
}

func (api *API) deleteShortUrl(w http.ResponseWriter, r *http.Request) error {
	// parse and validate input
	url, err := api.parseInput(r)
	if err != nil {
		return server.WriteError(w, *err)
	}

	err = api.svc.DeleteUrl(r.Context(), url)
	if err != nil {
		return server.WriteError(w, errors.NewErrorNotFound())
	}

	return server.Write(w, http.StatusOK, nil)
}

func (api *API) countRedirects(w http.ResponseWriter, r *http.Request) error {
	// parse and validate input
	url, err := api.parseInput(r)
	if err != nil {
		return server.WriteError(w, *err)
	}

	res, err := api.svc.CountRedirects(r.Context(), url)
	if err != nil {
		return server.WriteError(w, errors.NewErrorNotFound())
	}

	return server.Write(w, http.StatusOK, res)
}

func (api *API) redirect(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["shortId"]

	if id == "" {
		return server.WriteError(w, errors.NewErrorBadRequest())
	}

	url, err := api.svc.IncrementRedirect(r.Context(), id)
	if err != nil {
		return server.WriteError(w, errors.NewErrorNotFound())
	}

	http.Redirect(w, r, url, 301)
	return nil
}

func (api *API) parseInput(r *http.Request) (string, *errors.Error) {
	url := r.URL.Query().Get("url")

	// validate input
	url = strings.TrimSpace(url)
	if url == "" {
		api.le.Error("requested url is empty")
		e := errors.NewErrorBadRequest()
		return "", &e
	}

	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	return url, nil
}
