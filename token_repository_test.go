package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHttpTokenRepository(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	got := NewHttpTokenRepository(clock, mockSoundcloudApi{})
	assert.Equal(t, false, got.initialized)
}

func TestHttpTokenRepository_GetTokenOnFirstRunInitializes(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	got := NewHttpTokenRepository(clock, mockSoundcloudApi{})
	_, _ = got.GetToken()
	assert.Equal(t, true, got.initialized)
}

func TestHttpTokenRepository_GetTokenFirstRun(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	got, _ := NewHttpTokenRepository(clock, mockSoundcloudApi{}).GetToken()
	assert.Equal(t, "miao", got.AccessToken)
	assert.Equal(t, "bau", got.RefreshToken)
	assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
}

func TestHttpTokenRepository_GetTokenRenew(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	repository := HttpTokenRepository{
		currentToken: Token{ExpiresAt: time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)},
		sc:           mockSoundcloudApi{},
		initialized:  true,
		clock:        clock,
	}
	got, _ := repository.GetToken()
	assert.Equal(t, "miao_renewed", got.AccessToken)
	assert.Equal(t, "bau", got.RefreshToken)
	assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
}

func TestHttpTokenRepository_GetTokenAuthFailsError(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	_, err := NewHttpTokenRepository(clock, mockFailingSoundcloudApi{}).GetToken()
	assert.Equal(t, "auth_fail", err.Error())
}

func TestHttpTokenRepository_GetTokenRenewFailsError(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	repository := HttpTokenRepository{
		currentToken: Token{ExpiresAt: time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)},
		sc:           mockFailingSoundcloudApi{},
		initialized:  true,
		clock:        clock,
	}
	_, err := repository.GetToken()
	assert.Equal(t, "renew_fail", err.Error())
}

func TestHttpTokenRepository_GetTokenAuthJsonFailsError(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	_, err := NewHttpTokenRepository(clock, mockJsonFailingSoundcloudApi{}).GetToken()
	assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
}

func TestHttpTokenRepository_GetTokenRenewJsonFailsError(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	repository := HttpTokenRepository{
		currentToken: Token{ExpiresAt: time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)},
		sc:           mockJsonFailingSoundcloudApi{},
		initialized:  true,
		clock:        clock,
	}
	_, err := repository.GetToken()
	assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
}

type mockSoundcloudApi struct{}

func (m mockSoundcloudApi) Auth() ([]byte, error) {
	return []byte(`{"access_token":"miao","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`), nil
}

func (m mockSoundcloudApi) Renew(_ Token) ([]byte, error) {
	return []byte(`{"access_token":"miao_renewed","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`), nil
}

type mockFailingSoundcloudApi struct{}

func (m mockFailingSoundcloudApi) Auth() ([]byte, error) {
	return nil, errors.New("auth_fail")
}

func (m mockFailingSoundcloudApi) Renew(_ Token) ([]byte, error) {
	return nil, errors.New("renew_fail")
}

type mockJsonFailingSoundcloudApi struct{}

func (m mockJsonFailingSoundcloudApi) Auth() ([]byte, error) {
	return []byte(`12312`), nil
}

func (m mockJsonFailingSoundcloudApi) Renew(_ Token) ([]byte, error) {
	return []byte(`12312`), nil
}
