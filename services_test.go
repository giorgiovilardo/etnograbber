package main

import (
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpTrackDataService_GetTrackData(t *testing.T) {
	t.Run("should work ok, passes the token", func(t *testing.T) {
		expected := make(map[string]interface{})
		expected["id"] = 1
		expected["token"] = "bau"
		got, _ := NewHttpTrackDataService(mockTokenRepository{}, mockTrackDataRepository{}).GetTrackData(1)
		assert.Equal(t, got, expected)
	})

	t.Run("should emit token not available if token cannot be gained", func(t *testing.T) {
		_, err := NewHttpTrackDataService(newFailingMockTokenRepo("no token"), mockTrackDataRepository{}).GetTrackData(1)
		assert.Equal(t, err.Error(), "token not available")
	})

	t.Run("should emit track data not available if track data cannot be gained", func(t *testing.T) {
		_, err := NewHttpTrackDataService(mockTokenRepository{}, mockTrackDataRepository{wantErr: true, errMsg: "error"}).GetTrackData(1)
		assert.Equal(t, err.Error(), "trackData not available")
	})
}

func TestHttpCachedTrackService_GetTrack(t *testing.T) {
	t.Run("should work ok, pass the token", func(t *testing.T) {
		cache, _ := lru.New[int, []byte](1)
		got, _ := NewHttpCachedTrackService(cache, mockTokenRepository{}, mockTrackRepository{}).GetTrack(1)
		assert.Equal(t, got, []byte(`bau1`))
	})

	t.Run("should fetch from cache if available", func(t *testing.T) {
		cache := newMockLruCache()
		service := NewHttpCachedTrackService(cache, mockTokenRepository{}, mockTrackRepository{})
		_, _ = service.GetTrack(1)
		assert.False(t, cache.used)
		_, _ = service.GetTrack(1)
		assert.True(t, cache.used)
	})

	t.Run("should emit token not available if token cannot be gained", func(t *testing.T) {
		cache, _ := lru.New[int, []byte](1)
		_, err := NewHttpCachedTrackService(cache, newFailingMockTokenRepo("no token"), mockTrackRepository{}).GetTrack(1)
		assert.Equal(t, err.Error(), "token not available")
	})

	t.Run("should emit track data not available if track data cannot be gained", func(t *testing.T) {
		cache, _ := lru.New[int, []byte](1)
		_, err := NewHttpCachedTrackService(cache, mockTokenRepository{}, mockTrackRepository{wantErr: true, errMsg: "error"}).GetTrack(1)
		assert.Equal(t, err.Error(), "track not available")
	})
}

type mockTrackDataRepository struct {
	wantErr bool
	errMsg  string
}

func (m mockTrackDataRepository) GetTrackData(t Token, id int) (map[string]interface{}, error) {
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	track := make(map[string]interface{})
	track["id"] = id
	track["token"] = t.AccessToken
	return track, nil
}

type mockTokenRepository struct {
	wantErr bool
	errMsg  string
}

func (m mockTokenRepository) GetToken() (Token, error) {
	if m.wantErr {
		return Token{}, errors.New(m.errMsg)
	}
	return Token{AccessToken: "bau"}, nil
}

func newFailingMockTokenRepo(errMsg string) *mockTokenRepository {
	return &mockTokenRepository{wantErr: true, errMsg: errMsg}
}

type mockTrackRepository struct {
	wantErr bool
	errMsg  string
}

func (m mockTrackRepository) GetTrack(t Token, id int) ([]byte, error) {
	if m.wantErr {
		return nil, errors.New(m.errMsg)
	}
	return []byte(fmt.Sprintf("%s%d", t.AccessToken, id)), nil
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

func newMockLruCache() *mockLruCache {
	return &mockLruCache{cache: make(map[int][]byte)}
}
