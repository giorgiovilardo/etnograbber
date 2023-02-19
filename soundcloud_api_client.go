package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const AuthApiSuccessStatus = http.StatusOK

type HttpSoundcloudApi struct {
	c Config
}

func NewHttpSoundcloudApi(c Config) *HttpSoundcloudApi {
	return &HttpSoundcloudApi{c: c}
}

func (s *HttpSoundcloudApi) GetTrackData(t Token, id int) (map[string]interface{}, error) {
	trackUrl := fmt.Sprintf("%s/%d", s.c.BaseApiUrl, id)
	authHeader := fmt.Sprintf("OAuth %s", t.AccessToken)
	client := http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest(http.MethodGet, trackUrl, nil)
	req.Header.Set("Authorization", authHeader)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get track data"), err)
	}

	result := make(map[string]interface{})

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse track data"), err)
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, errors.Join(errors.New("failed to jsonize track data"), err)
	}

	return result, nil
}

func (s *HttpSoundcloudApi) GetTrack(t Token, id int) ([]byte, error) {
	trackUrl := fmt.Sprintf("%s/%d/stream", s.c.BaseApiUrl, id)
	authHeader := fmt.Sprintf("OAuth %s", t.AccessToken)
	client := http.Client{Timeout: time.Second * 20}
	req, _ := http.NewRequest(http.MethodGet, trackUrl, nil)
	req.Header.Set("Authorization", authHeader)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get track stream"), err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get track stream"), err)
	}

	return body, nil
}

func (s *HttpSoundcloudApi) Auth() ([]byte, error) {
	tokenData, err := s.getToken()
	if err == nil {
		return tokenData, nil
	}

	fallbackTokenData, fallbackErr := s.getFallback()
	if fallbackErr != nil {
		return nil, errors.Join(errors.New("impossible to acquire a token"), fallbackErr, err)
	}

	return fallbackTokenData, nil
}

func (s *HttpSoundcloudApi) Renew(t Token) ([]byte, error) {
	formData := url.Values{}
	formData.Add("grant_type", "refresh_token")
	formData.Add("client_id", s.c.ClientId)
	formData.Add("client_secret", s.c.ClientSecret)
	formData.Add("refresh_token", t.RefreshToken)

	client := http.Client{Timeout: time.Second * 5}

	res, err := client.Post(s.c.BaseAuthUrl, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, errors.Join(errors.New("could not renew the token, post failed"), err)
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("unparsable api body"), err)
	}

	return result, nil
}

func (s *HttpSoundcloudApi) getToken() ([]byte, error) {
	formData := url.Values{}
	formData.Add("grant_type", "client_credentials")
	formData.Add("client_id", s.c.ClientId)
	formData.Add("client_secret", s.c.ClientSecret)

	client := http.Client{Timeout: time.Second * 5}

	res, err := client.Post(s.c.BaseAuthUrl, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, errors.Join(errors.New("failed to get token from BaseAuth, network error"), err)
	}

	if res.StatusCode != AuthApiSuccessStatus {
		return nil, errors.New(fmt.Sprintf("failed to get token from BaseAuth, status not %d", AuthApiSuccessStatus))
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("unparsable api body"), err)
	}

	return result, nil
}

func (s *HttpSoundcloudApi) getFallback() ([]byte, error) {
	client := http.Client{Timeout: time.Second * 5}

	res, err := client.Get(s.c.FallbackAuthUrl)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get token from FallbackAuth"), err)
	}

	if res.StatusCode != AuthApiSuccessStatus {
		return nil, errors.New(fmt.Sprintf("failed to get token from FallbackAuth, status not %d", AuthApiSuccessStatus))
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("unparsable api body"), err)
	}

	return result, nil
}
