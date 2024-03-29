package main

import (
	"errors"
	clockLib "github.com/giorgiovilardo/etnograbber/clock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestHttpTokenRepository_GetToken(t *testing.T) {
	clock := clockLib.NewBrokenClock(time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC))
	repoWithOldToken := func(api SoundcloudApi) *HttpTokenRepository {
		return &HttpTokenRepository{
			currentToken: Token{ExpiresAt: time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)},
			sc:           api,
			initialized:  true,
			clock:        clock,
		}
	}

	t.Run("should return a token", func(t *testing.T) {
		got, _ := NewHttpTokenRepository(clock, &mockSoundcloudApi{}).GetToken()
		assert.Equal(t, "miao", got.AccessToken)
		assert.Equal(t, "bau", got.RefreshToken)
		assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
	})

	t.Run("should transparently renew the token if expired", func(t *testing.T) {
		got, _ := repoWithOldToken(&mockSoundcloudApi{}).GetToken()
		assert.Equal(t, "miao_renewed", got.AccessToken)
		assert.Equal(t, "bau", got.RefreshToken)
		assert.Equal(t, clock.Now().Add(time.Second*time.Duration(got.ExpiresIn)), got.ExpiresAt)
	})

	t.Run("should return the same error from the api if auth fails", func(t *testing.T) {
		_, err := NewHttpTokenRepository(clock, &mockSoundcloudApi{wantErr: true, errMsg: "auth_fail"}).GetToken()
		assert.Equal(t, "auth_fail", err.Error())
	})

	t.Run("should return the same error from the api if renew fails", func(t *testing.T) {
		_, err := repoWithOldToken(&mockSoundcloudApi{wantErr: true, errMsg: "renew_fail"}).GetToken()
		assert.Equal(t, "renew_fail", err.Error())
	})

	t.Run("should fail if auth can't deserialize json into token", func(t *testing.T) {
		_, err := NewHttpTokenRepository(clock, &mockSoundcloudApi{thatReturns: []byte(`12345`)}).GetToken()
		assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
	})

	t.Run("should fail if renew can't deserialize json into token", func(t *testing.T) {
		_, err := repoWithOldToken(&mockSoundcloudApi{thatReturns: []byte(`12345`)}).GetToken()
		assert.Equal(t, "json: cannot unmarshal number into Go value of type main.Token", err.Error())
	})

	t.Run("should not auth more than once even under concurrent access", func(t *testing.T) {
		wg := sync.WaitGroup{}
		tries := 20
		api := &mockSoundcloudApi{}
		repo := NewHttpTokenRepository(clock, api)
		wg.Add(tries)
		for i := 0; i < tries; i++ {
			go func() {
				_, _ = repo.GetToken()
				wg.Done()
			}()
		}
		wg.Wait()
		assert.Equal(t, 1, api.Calls)
	})

	t.Run("should not renew more than once even under concurrent access", func(t *testing.T) {
		wg := sync.WaitGroup{}
		tries := 20
		api := &mockSoundcloudApi{}
		repo := repoWithOldToken(api)
		wg.Add(tries)
		for i := 0; i < tries; i++ {
			go func() {
				_, _ = repo.GetToken()
				wg.Done()
			}()
		}
		wg.Wait()
		assert.Equal(t, 1, api.Calls)
	})
}

type mockSoundcloudApi struct {
	thatReturns []byte
	wantErr     bool
	errMsg      string
	Calls       int
}

func (m *mockSoundcloudApi) Auth() ([]byte, error) {
	m.Calls += 1
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	if m.thatReturns != nil {
		return m.thatReturns, nil
	}
	return []byte(`{"access_token":"miao","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`), nil
}

func (m *mockSoundcloudApi) Renew(_ Token) ([]byte, error) {
	m.Calls += 1
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	if m.thatReturns != nil {
		return m.thatReturns, nil
	}
	return []byte(`{"access_token":"miao_renewed","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`), nil
}
