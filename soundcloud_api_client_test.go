package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpSoundcloudApi_Renew(t *testing.T) {
	t.Run("should return the raw response", func(t *testing.T) {
		expected := []byte(`{"micio":"miao"}`)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write(expected)
		}))
		defer server.Close()
		conf := Config{BaseAuthUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.Renew(Token{})
		assert.Equal(t, expected, res)
	})

	t.Run("should post the correct data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "refresh_token", r.PostFormValue("grant_type"))
			assert.Equal(t, "micio", r.PostFormValue("client_id"))
			assert.Equal(t, "mao", r.PostFormValue("client_secret"))
			assert.Equal(t, "reftoken", r.PostFormValue("refresh_token"))
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(AuthApiSuccessStatus)
		}))
		defer server.Close()
		conf := Config{BaseAuthUrl: server.URL, ClientId: "micio", ClientSecret: "mao"}
		api := NewHttpSoundcloudApi(conf)
		_, _ = api.Renew(Token{RefreshToken: "reftoken"})
	})

	t.Run("should fail if it can't make the request", func(t *testing.T) {
		conf := Config{BaseAuthUrl: "bad server"}
		api := NewHttpSoundcloudApi(conf)
		res, err := api.Renew(Token{})
		assert.Equal(t, []byte(nil), res)
		assert.Contains(t, err.Error(), "could not renew the token, post failed")
	})
}

func TestHttpSoundcloudApi_Auth(t *testing.T) {
	t.Run("should fetch the token from base auth url", func(t *testing.T) {
		expected := []byte(`{"ola":"ola"}`)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write(expected)
		}))
		defer server.Close()
		conf := Config{BaseAuthUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.Auth()
		assert.Equal(t, expected, res)
	})

	t.Run("should fetch the token from fallback url if base has network error", func(t *testing.T) {
		expected := []byte(`{"ola":"ola"}`)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write(expected)
		}))
		defer server.Close()
		conf := Config{BaseAuthUrl: "not a server", FallbackAuthUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.Auth()
		assert.Equal(t, expected, res)
	})

	t.Run("should use the fallback auth if base has server errors", func(t *testing.T) {
		expected := []byte(`{"ola":"ola"}`)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		fallbackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write(expected)
		}))
		defer server.Close()
		defer fallbackServer.Close()
		conf := Config{BaseAuthUrl: server.URL, FallbackAuthUrl: fallbackServer.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.Auth()
		assert.Equal(t, expected, res)
	})

	t.Run("should error with both base and fallback auth unavailable", func(t *testing.T) {
		conf := Config{BaseAuthUrl: "bad server", FallbackAuthUrl: "bad server"}
		api := NewHttpSoundcloudApi(conf)
		res, err := api.Auth()
		assert.Equal(t, []byte(nil), res)
		assert.Contains(t, err.Error(), "impossible to acquire a token")
	})

	t.Run("should post client auth data to the base auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "client_credentials", r.PostFormValue("grant_type"))
			assert.Equal(t, "micio", r.PostFormValue("client_id"))
			assert.Equal(t, "mao", r.PostFormValue("client_secret"))
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(AuthApiSuccessStatus)
		}))
		defer server.Close()
		conf := Config{BaseAuthUrl: server.URL, ClientId: "micio", ClientSecret: "mao"}
		api := NewHttpSoundcloudApi(conf)
		_, _ = api.Auth()
	})

	t.Run("should get with no data sent to fallback auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NotEqual(t, "client_credentials", r.PostFormValue("grant_type"))
			assert.NotEqual(t, "micio", r.PostFormValue("client_id"))
			assert.NotEqual(t, "mao", r.PostFormValue("client_secret"))
			assert.NotEqual(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(AuthApiSuccessStatus)
		}))
		defer server.Close()
		conf := Config{FallbackAuthUrl: server.URL, ClientId: "micio", ClientSecret: "mao"}
		api := NewHttpSoundcloudApi(conf)
		_, _ = api.Auth()
	})

	t.Run("should fail if both auths have server errors", func(t *testing.T) {
		expected := []byte(`{"ola":"ola"}`)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		fallbackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(expected)
		}))
		defer server.Close()
		defer fallbackServer.Close()
		conf := Config{BaseAuthUrl: server.URL, FallbackAuthUrl: fallbackServer.URL}
		api := NewHttpSoundcloudApi(conf)
		res, err := api.Auth()
		assert.Equal(t, []byte(nil), res)
		assert.Contains(t, err.Error(), "failed to get token from BaseAuth, status not 200")
		assert.Contains(t, err.Error(), "failed to get token from FallbackAuth, status not 200")
	})
}

func TestHttpSoundcloudApi_GetTrackData(t *testing.T) {
	t.Run("should work correctly", func(t *testing.T) {
		expected := map[string]interface{}{"fake": "result"}
		token := Token{AccessToken: "faketoken"}
		id := 100
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "OAuth faketoken", r.Header.Get("Authorization"))
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, fmt.Sprintf("/%d", id), r.URL.RequestURI())
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write([]byte(`{"fake":"result"}`))
		}))
		defer server.Close()
		conf := Config{BaseApiUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.GetTrackData(token, id)
		assert.Equal(t, expected, res)
	})

	t.Run("should fail with network error", func(t *testing.T) {
		conf := Config{BaseApiUrl: "YOLO"}
		api := NewHttpSoundcloudApi(conf)
		_, err := api.GetTrackData(Token{}, 0)
		assert.Contains(t, err.Error(), "failed to get track data")
	})

	t.Run("should fail with non 200 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		conf := Config{BaseApiUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		_, err := api.GetTrackData(Token{}, 0)
		assert.Contains(t, err.Error(), "failed to get track data")
	})
}

func TestHttpSoundcloudApi_GetTrack(t *testing.T) {
	t.Run("should work correctly", func(t *testing.T) {
		token := Token{AccessToken: "faketoken"}
		id := 100
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "OAuth faketoken", r.Header.Get("Authorization"))
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, fmt.Sprintf("/%d/stream", id), r.URL.RequestURI())
			w.WriteHeader(AuthApiSuccessStatus)
			_, _ = w.Write([]byte(`{"fake":"result"}`))
		}))
		defer server.Close()
		conf := Config{BaseApiUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		res, _ := api.GetTrack(token, id)
		assert.Equal(t, []byte(`{"fake":"result"}`), res)
	})

	t.Run("should fail with network error", func(t *testing.T) {
		conf := Config{BaseApiUrl: "YOLO"}
		api := NewHttpSoundcloudApi(conf)
		_, err := api.GetTrack(Token{}, 0)
		assert.Contains(t, err.Error(), "failed to get track stream")
	})

	t.Run("should fail with non 200 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		conf := Config{BaseApiUrl: server.URL}
		api := NewHttpSoundcloudApi(conf)
		_, err := api.GetTrack(Token{}, 0)
		assert.Contains(t, err.Error(), "failed to get track stream")
	})
}
