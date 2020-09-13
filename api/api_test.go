package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/shortener"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testUrl  = "http://www.test.com"
	shortUrl = "http://www.short.me/RMAp1Vz"
	shortId  = "RMAp1Vz"
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

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api?url=%s", testUrl), nil)

	svc.EXPECT().ShortenUrl(req.Context(), testUrl).Return(shortId, nil)

	resp := httptest.NewRecorder()
	if err := api.createShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	decoder := json.NewDecoder(resp.Body)
	var payload ResponseApi
	require.NoError(t, decoder.Decode(&payload))

	// example.com is the host used for testing
	require.Equal(t, "ok", payload.Status)
	require.Equal(t, fmt.Sprintf("http://example.com/%s", shortId), payload.Url)
	require.Equal(t, "create", payload.Operation)
}

func TestAPI_ReadShortUrl(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api?url=%s", shortUrl), nil)

	svc.EXPECT().RetrieveUrl(req.Context(), shortUrl).Return(shortId, nil)

	resp := httptest.NewRecorder()
	if err := api.readShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	decoder := json.NewDecoder(resp.Body)
	var payload ResponseApi
	require.NoError(t, decoder.Decode(&payload))

	require.Equal(t, "ok", payload.Status)
	require.Equal(t, "read", payload.Operation)
	require.Equal(t, shortId, payload.Url)
}

func TestAPI_DeleteShortUrl(t *testing.T) {
	api, svc := MakeTestApi(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api?url=%s", shortUrl), nil)

	svc.EXPECT().DeleteUrl(req.Context(), shortUrl)

	resp := httptest.NewRecorder()
	if err := api.deleteShortUrl(resp, req); err != nil {
		t.Error(err)
	}

	verifyStatus(t, http.StatusOK, resp.Code)

	decoder := json.NewDecoder(resp.Body)
	var payload ResponseApi
	require.NoError(t, decoder.Decode(&payload))

	require.Equal(t, "ok", payload.Status)
	require.Equal(t, "delete", payload.Operation)
	require.Equal(t, shortUrl, payload.Url)
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

	decoder := json.NewDecoder(resp.Body)
	var payload ResponseCount
	require.NoError(t, decoder.Decode(&payload))

	require.Equal(t, "ok", payload.Status)
	require.Equal(t, "count", payload.Operation)
	require.Equal(t, int64(2), payload.Count)
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

	req := httptest.NewRequest(http.MethodPost, "/api?url=", nil)

	resp := httptest.NewRecorder()
	_ = api.createShortUrl(resp, req)

	verifyStatus(t, http.StatusBadRequest, resp.Code)
}

func verifyStatus(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("Status code was %d, expected %d ", actual, expected)
	}
}
