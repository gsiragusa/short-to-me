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
			Path:    "/api",
			Handler: api.createShortUrl,
		},
		{
			Name:    "read-short-url",
			Method:  http.MethodGet,
			Path:    "/api",
			Handler: api.readShortUrl,
		},
		{
			Name:    "delete-short-url",
			Method:  http.MethodDelete,
			Path:    "/api",
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

// Parses the input url and generates a short one for it
func (api *API) createShortUrl(w http.ResponseWriter, r *http.Request) error {
	// swagger:operation POST /api Api createShortUrl
	// Shorten url
	//
	// Consumes a url and shortens it
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: url
	//   in: query
	//   description: url to shorten
	//   required: true
	//   type: string
	//
	// responses:
	//   '200':
	//     description: "Short url"
	//     schema:
	//       type: object
	//       properties:
	//         url:
	//           type: string
	//           example: "http://www.example.com/RMAp1Vz"
	//   '400':
	//     description: Not Found
	//   '404':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

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
	// swagger:operation GET /api Api readShortUrl
	// Read short url
	//
	// Receives a short url and returns the extended one
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: url
	//   in: query
	//   description: short url to read
	//   required: true
	//   type: string
	//
	// responses:
	//   '200':
	//     description: "Extended url"
	//     schema:
	//       type: object
	//       properties:
	//         url:
	//           type: string
	//           example: "https://www.google.com"
	//   '400':
	//     description: Not Found
	//   '404':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

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
	// swagger:operation DELETE /api Api deleteShortUrl
	// Delete short url
	//
	// Deletes the short url
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: url
	//   in: query
	//   description: short url to delete
	//   required: true
	//   type: string
	//
	// responses:
	//   '200':
	//     description: "Operation result"
	//     schema:
	//       type: object
	//       properties:
	//         result:
	//           type: string
	//           example: "ok"
	//   '400':
	//     description: Not Found
	//   '404':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

	// parse and validate input
	url, err := api.parseInput(r)
	if err != nil {
		return server.WriteError(w, *err)
	}

	err = api.svc.DeleteUrl(r.Context(), url)
	if err != nil {
		return server.WriteError(w, errors.NewErrorNotFound())
	}

	return server.Write(w, http.StatusOK, &shortener.Operation{Result: "ok"})
}

func (api *API) countRedirects(w http.ResponseWriter, r *http.Request) error {
	// swagger:operation GET /api/count Api countRedirects
	// Count number of redirects
	//
	// Returns the number of redirections for a given short url
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: url
	//   in: query
	//   description: short url to count redirections
	//   required: true
	//   type: string
	//
	// responses:
	//   '200':
	//     description: "Total count"
	//     schema:
	//       type: integer
	//       format: int64
	//   '400':
	//     description: Not Found
	//   '404':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

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
	// swagger:operation GET /{shortId} Redirect countRedirects
	// Redirect to extended url
	//
	// Redirect to extended url
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: shortId
	//   in: path
	//   description: the id of the short url for redirection
	//   required: true
	//   type: string
	//
	// responses:
	//   '301':
	//     description: "Redirects to extended url"
	//   '400':
	//     description: Not Found
	//   '404':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

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
