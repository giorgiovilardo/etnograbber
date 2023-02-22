package main

import (
	"errors"
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
		if assert.NoError(t, TrackDataHandler(mockTrackDataService{})(c)) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if path param is non int castable", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackId not a number"}`
		c, r := setupEcho("aba")
		if assert.NoError(t, TrackDataHandler(mockTrackDataService{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't acquire a token", func(t *testing.T) {
		expectedResponseBody := `{"error":"token not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackDataHandler(newFailingMockTrackDataService("token not available"))(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't read track data", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackData not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackDataHandler(newFailingMockTrackDataService("trackData not available"))(c)) {
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
		if assert.NoError(t, TrackHandler(mockTrackService{})(c)) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "audio/mpeg", r.Header().Get("Content-Type"))
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if path param is non int castable", func(t *testing.T) {
		expectedResponseBody := `{"error":"trackId not a number"}`
		c, r := setupEcho("aba")
		if assert.NoError(t, TrackHandler(mockTrackService{})(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't acquire a token", func(t *testing.T) {
		expectedResponseBody := `{"error":"token not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackHandler(newFailingMockTrackService("token not available"))(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})

	t.Run("should fail if it can't read track stream", func(t *testing.T) {
		expectedResponseBody := `{"error":"track not available"}`
		c, r := setupEcho("1234")
		if assert.NoError(t, TrackHandler(newFailingMockTrackService("track not available"))(c)) {
			assert.Equal(t, http.StatusServiceUnavailable, r.Code)
			assert.Equal(t, expectedResponseBody, strings.Trim(r.Body.String(), "\n"))
		}
	})
}

type mockTrackDataService struct {
	wantErr bool
	errMsg  string
}

func (m mockTrackDataService) GetTrackData(id int) (map[string]interface{}, error) {
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	track := make(map[string]interface{})
	track["id"] = id
	return track, nil
}

func newFailingMockTrackDataService(errMsg string) *mockTrackDataService {
	return &mockTrackDataService{wantErr: true, errMsg: errMsg}
}

type mockTrackService struct {
	wantErr bool
	errMsg  string
}

func (m mockTrackService) GetTrack(_ int) ([]byte, error) {
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	return []byte(`yolo`), nil
}

func newFailingMockTrackService(errMsg string) *mockTrackService {
	return &mockTrackService{wantErr: true, errMsg: errMsg}
}
