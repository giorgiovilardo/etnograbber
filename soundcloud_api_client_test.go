package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpSoundcloudApi_RenewReturnsByteArrayOfTheResponse(t *testing.T) {
	expected := []byte(`{"micio":"miao"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(expected)
	}))
	defer server.Close()
	conf := Config{BaseAuthUrl: server.URL}
	api := NewHttpSoundcloudApi(conf)
	res, _ := api.Renew(Token{})
	assert.Equal(t, expected, res)
}

func TestHttpSoundcloudApi_RenewPostsTheCorrectData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "refresh_token", r.PostFormValue("grant_type"))
		assert.Equal(t, "micio", r.PostFormValue("client_id"))
		assert.Equal(t, "mao", r.PostFormValue("client_secret"))
		assert.Equal(t, "reftoken", r.PostFormValue("refresh_token"))
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		w.WriteHeader(200)
	}))
	conf := Config{BaseAuthUrl: server.URL, ClientId: "micio", ClientSecret: "mao"}
	defer server.Close()
	api := NewHttpSoundcloudApi(conf)
	_, _ = api.Renew(Token{RefreshToken: "reftoken"})
}

func TestHttpSoundcloudApi_RenewErrors(t *testing.T) {
	conf := Config{BaseAuthUrl: "bad server"}
	api := NewHttpSoundcloudApi(conf)
	res, err := api.Renew(Token{})
	assert.Equal(t, []byte(nil), res)
	assert.Contains(t, err.Error(), "could not renew the token, post failed")
}
