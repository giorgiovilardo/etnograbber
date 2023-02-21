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
	t.Run("should return a small message signaling service is alive", func(t *testing.T) {
		expectedResponseBody := `{"message":"ok"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, HealthHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		}
	})
}

func TestTrackDataHandler(t *testing.T) {
	setupEcho := func(trackId string) (echo.Context, *httptest.ResponseRecorder) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/"+trackId, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:trackId")
		c.SetParamNames("trackId")
		c.SetParamValues(trackId)
		return c, rec
	}

	t.Run("should proxy over what the repo returns", func(t *testing.T) {
		expectedResponseBody := `{"id":1234}` // change here to see me fail
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackDataHandler(mockTokenRepository{}, mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if path param is non int castable", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackId not a number"}`
		c, r := setupEcho("aba")
		if assert.NoError(t, TrackDataHandler(mockTokenRepository{}, mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't acquire a token", func(t *testing.T) {
		expectedResponseBody := `{"error":"token not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackDataHandler(newFailingMockTokenRepo("i failed"), mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't read track data", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackData not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackDataHandler(mockTokenRepository{}, newFailingMockTrackRepository("i failed"))(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})
}

func TestTrackHandler(t *testing.T) {
	setupEcho := func(trackId string) (echo.Context, *httptest.ResponseRecorder) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/"+trackId+"/stream", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:trackId")
		c.SetParamNames("trackId")
		c.SetParamValues(trackId)
		return c, rec
	}

	t.Run("should proxy over what the repo returns, with the audio/mpeg content-type", func(t *testing.T) {
		expectedResponseBody := `yolo`
		c, r := setupEcho("1234")
		cache, _ := lru.New[int, []byte](1)
		if assert.NoError(t, TrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "audio/mpeg", r.Header().Get("Content-Type"))
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if path param is non int castable", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackId not a number"}`
		c, r := setupEcho("aba")
		cache, _ := lru.New[int, []byte](1)
		if assert.NoError(t, TrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't acquire a token", func(t *testing.T) {
		expectedResponseBody := `{"error":"token not available"}`
		c, r := setupEcho("1234")
		cache, _ := lru.New[int, []byte](1)
		if assert.NoError(t, TrackHandler(cache, newFailingMockTokenRepo("i failed"), mockTrackRepository{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't read track stream", func(t *testing.T) {
		expectedResponseBody := `{"error":"track not available"}`
		c, r := setupEcho("1234")
		cache, _ := lru.New[int, []byte](1)
		if assert.NoError(t, TrackHandler(cache, mockTokenRepository{}, newFailingMockTrackRepository("i failed"))(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})
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
	handler := TrackHandler(cache, mockTokenRepository{}, mockTrackRepository{})
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

type mockTokenRepository struct {
	withErr bool
	errMsg  string
}

func (m mockTokenRepository) GetToken() (Token, error) {
	if m.withErr {
		return Token{}, errors.New(m.errMsg)
	}
	return Token{AccessToken: "bau"}, nil
}

func newFailingMockTokenRepo(errMsg string) *mockTokenRepository {
	return &mockTokenRepository{withErr: true, errMsg: errMsg}
}

type mockTrackRepository struct {
	withErr bool
	errMsg  string
}

func (m mockTrackRepository) GetTrackData(_ Token, id int) (map[string]interface{}, error) {
	if m.withErr {
		return nil, errors.New(m.errMsg)
	}
	track := make(map[string]interface{})
	track["id"] = id
	return track, nil
}

func (m mockTrackRepository) GetTrack(_ Token, _ int) ([]byte, error) {
	if m.withErr {
		return nil, errors.New(m.errMsg)
	}
	return []byte(`yolo`), nil
}

func newFailingMockTrackRepository(errMsg string) *mockTrackRepository {
	return &mockTrackRepository{withErr: true, errMsg: errMsg}
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
