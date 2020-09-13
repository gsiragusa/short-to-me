package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/shortener"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testUrl = "http://www.test.com"
	shortUrl = "http://www.short.me/RMAp1Vz"
	shortId = "RMAp1Vz"
)

func MakeTestApi(t *testing.T) (*API, *shortener.MockService) {
	log := logrus.New()
	log.Out = ioutil.Discard // silent logger

	conf, err := config.Configure()
	require.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := shortener.NewMockService(ctrl)

	return NewAPI(log, conf, svc), svc
}

func TestAPI_CreateShortUrl(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/shortener?url=%s", testUrl), nil)

	svc.EXPECT().ShortenUrl(req.Context(), testUrl).Return("123", nil)

	resp := httptest.NewRecorder()
	if err := api.createShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	decoder := json.NewDecoder(resp.Body)
	var payload shortener.ModelShorten
	require.NoError(t, decoder.Decode(&payload))

	// example.com is the host used for testing
	require.Equal(t, "http://example.com/123", payload.Url)
}

func TestAPI_ReadShortUrl(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shortener?url=%s", shortUrl), nil)

	svc.EXPECT().RetrieveUrl(req.Context(), shortUrl).Return(shortId, nil)

	resp := httptest.NewRecorder()
	if err := api.readShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	decoder := json.NewDecoder(resp.Body)
	var payload shortener.ModelShorten
	require.NoError(t, decoder.Decode(&payload))

	require.Equal(t, shortId, payload.Url)
}

func TestAPI_DeleteShortUrl(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/shortener?url=%s", shortUrl), nil)

	svc.EXPECT().DeleteUrl(req.Context(), shortUrl)

	resp := httptest.NewRecorder()
	if err := api.deleteShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)
}

func TestAPI_CountRedirects(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/count?url=%s", shortUrl), nil)

	svc.EXPECT().CountRedirects(req.Context(), shortUrl).Return(int64(2), nil)

	resp := httptest.NewRecorder()
	if err := api.countRedirects(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	count, err := strconv.Atoi(resp.Body.String())
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestAPI_Redirect(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/123"), nil)
	req = mux.SetURLVars(req, map[string]string{"shortId": "123"})

	svc.EXPECT().IncrementRedirect(req.Context(), "123")

	resp := httptest.NewRecorder()
	if err := api.redirect(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, 301, resp.Code)
}

func TestAPI_BadRequest(t *testing.T) {
	api, _ := MakeTestApi(t)

	req:= httptest.NewRequest(http.MethodPost, "/api/shortener?url=", nil)

	resp := httptest.NewRecorder()
	_ = api.createShortUrl(resp, req)

	verifyStatus(t, http.StatusBadRequest, resp.Code)
}

func verifyStatus(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("Status code was %d, expected %d ", actual, expected)
	}
}
