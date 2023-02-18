package main

import (
	"bytes"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	expectedResponseBody := `{"message":"ok"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, HealthHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestTrackHandlerHasAPathParam(t *testing.T) {
	expectedResponseBody := `{"id":1234}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, TrackHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestTrackHandlerFailsIfPathParamIsNotIntParsable(t *testing.T) {
	expectedResponseBody := `{"error":"trackId not a number"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/aba", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("aba")
	if assert.NoError(t, TrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestTrackHandlerFailsIfTokenNotAvailable(t *testing.T) {
	expectedResponseBody := `{"error":"token not available"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockFailingTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, TrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestTrackHandlerFailsIfTrackNotAvailable(t *testing.T) {
	expectedResponseBody := `{"error":"trackData not available"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockFailingTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, TrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandler(t *testing.T) {
	expectedResponseBody := `yolo`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, StreamTrackHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "audio/mpeg", rec.Header().Get("Content-Type"))
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandlerFailsIfPathParamIsNotIntParsable(t *testing.T) {
	expectedResponseBody := `{"error":"trackId not a number"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/aba/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("aba")
	if assert.NoError(t, StreamTrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandlerFailsIfTokenNotAvailable(t *testing.T) {
	expectedResponseBody := `{"error":"token not available"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockTrackRepository{})
	c.Set("tokenRepository", mockFailingTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, StreamTrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandlerFailsIfTrackNotAvailable(t *testing.T) {
	expectedResponseBody := `{"error":"trackData not available"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("trackRepository", mockFailingTrackRepository{})
	c.Set("tokenRepository", mockTokenRepository{})
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	if assert.NoError(t, StreamTrackHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

type mockTokenRepository struct{}

func (m mockTokenRepository) GetToken() (Token, error) {
	return Token{AccessToken: "bau"}, nil
}

type mockTrackRepository struct{}

func (m mockTrackRepository) GetTrackData(t Token, id int) (map[string]interface{}, error) {
	track := make(map[string]interface{})
	track["id"] = id
	return track, nil
}

func (m mockTrackRepository) GetTrack(t Token, id int) (io.Reader, error) {
	return bytes.NewReader([]byte(`yolo`)), nil
}

type mockFailingTokenRepository struct{}

func (m mockFailingTokenRepository) GetToken() (Token, error) {
	return Token{}, errors.New("i failed")
}

type mockFailingTrackRepository struct{}

func (m mockFailingTrackRepository) GetTrackData(t Token, id int) (map[string]interface{}, error) {
	return nil, errors.New("random error")
}

func (m mockFailingTrackRepository) GetTrack(t Token, id int) (io.Reader, error) {
	return nil, errors.New("random error")
}
