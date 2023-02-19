package main

import (
	"errors"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	cache, _ := lru.New[int, []byte](1)
	if assert.NoError(t, StreamTrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "audio/mpeg", rec.Header().Get("Content-Type"))
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandlerFetchesFromCacheIfAvailable(t *testing.T) {
	expectedResponseBody := `yolo`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:trackId/stream")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	cache := &mockLruCache{
		cache: map[int][]byte{},
		used:  false,
	}
	handler := StreamTrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "audio/mpeg", rec.Header().Get("Content-Type"))
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		assert.False(t, cache.used)
	}
	if assert.NoError(t, handler(c)) {
		assert.True(t, cache.used)
	}
}

func TestStreamTrackHandlerFailsIfPathParamIsNotIntParsable(t *testing.T) {
	expectedResponseBody := `{"error":"trackId not a number"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/aba/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("aba")
	cache, _ := lru.New[int, []byte](1)
	if assert.NoError(t, StreamTrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})(c)) {
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
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	cache, _ := lru.New[int, []byte](1)
	if assert.NoError(t, StreamTrackHandler(cache, mockFailingTokenRepository{}, mockTrackRepository{})(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestStreamTrackHandlerFailsIfTrackNotAvailable(t *testing.T) {
	expectedResponseBody := `{"error":"track not available"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/1234/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:trackId")
	c.SetParamNames("trackId")
	c.SetParamValues("1234")
	cache, _ := lru.New[int, []byte](1)
	if assert.NoError(t, StreamTrackHandler(cache, mockTokenRepository{}, mockFailingTrackRepository{})(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}

type mockTokenRepository struct{}

func (m mockTokenRepository) GetToken() (Token, error) {
	return Token{AccessToken: "bau"}, nil
}

type mockTrackRepository struct{}

func (m mockTrackRepository) GetTrackData(_ Token, id int) (map[string]interface{}, error) {
	track := make(map[string]interface{})
	track["id"] = id
	return track, nil
}

func (m mockTrackRepository) GetTrack(_ Token, _ int) ([]byte, error) {
	return []byte(`yolo`), nil
}

type mockFailingTokenRepository struct{}

func (m mockFailingTokenRepository) GetToken() (Token, error) {
	return Token{}, errors.New("i failed")
}

type mockFailingTrackRepository struct{}

func (m mockFailingTrackRepository) GetTrackData(_ Token, _ int) (map[string]interface{}, error) {
	return nil, errors.New("random error")
}

func (m mockFailingTrackRepository) GetTrack(_ Token, _ int) ([]byte, error) {
	return nil, errors.New("random error")
}

type mockLruCache struct {
	cache map[int][]byte
	used  bool
}

func (m *mockLruCache) Add(key int, value []byte) (evicted bool) {
	m.cache[key] = value
	return false
}

func (m *mockLruCache) Contains(key int) bool {
	_, ok := m.cache[key]
	return ok
}

func (m *mockLruCache) Get(key int) (value []byte, ok bool) {
	m.used = true
	return m.cache[key], true
}
