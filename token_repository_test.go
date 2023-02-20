package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHttpTokenRepository_GetToken(t *testing.T) {
	clock := NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	repoWith := func(api SoundcloudApi) *HttpTokenRepository {
		return &HttpTokenRepository{
			currentToken: Token{ExpiresAt: time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)},
			sc:           api,
			initialized:  true,
			clock:        clock,
		}
	}

	t.Run("should return a token", func(t *testing.T) {
		got, _ := NewHttpTokenRepository(clock, mockSoundcloudApi{}).GetToken()
		assert.Equal(t, "miao", got.AccessToken)
		assert.Equal(t, "bau", got.RefreshToken)
		assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
	})

	t.Run("should transparently renew the token if expired", func(t *testing.T) {
		got, _ := repoWith(mockSoundcloudApi{}).GetToken()
		assert.Equal(t, "miao_renewed", got.AccessToken)
		assert.Equal(t, "bau", got.RefreshToken)
		assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
	})

	t.Run("should return the same error from the api if auth fails", func(t *testing.T) {
		_, err := NewHttpTokenRepository(clock, mockFailingSoundcloudApi{}).GetToken()
		assert.Equal(t, "auth_fail", err.Error())
	})

	t.Run("should return the same error from the api if renew fails", func(t *testing.T) {
		_, err := repoWith(mockFailingSoundcloudApi{}).GetToken()
		assert.Equal(t, "renew_fail", err.Error())
	})

	t.Run("should fail if auth can't deserialize json into token", func(t *testing.T) {
		_, err := NewHttpTokenRepository(clock, mockJsonFailingSoundcloudApi{}).GetToken()
		assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
	})

	t.Run("should fail if renew can't deserialize json into token", func(t *testing.T) {
		_, err := repoWith(mockJsonFailingSoundcloudApi{}).GetToken()
		assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
	})
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
