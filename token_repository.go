package main

import (
	"github.com/giorgiovilardo/etnograbber/clock"
	"sync"
	"time"
)

type HttpTokenRepository struct {
	currentToken Token
	sc           SoundcloudApi
	initialized  bool
	clock        clock.Clock
	Mu           sync.Mutex
}

func NewHttpTokenRepository(c clock.Clock, sc SoundcloudApi) *HttpTokenRepository {
	return &HttpTokenRepository{
		initialized: false,
		clock:       c,
		sc:          sc,
	}
}

func (s *HttpTokenRepository) GetToken() (Token, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if !s.initialized {
		token, err := newToken(s.sc, s.clock.Now())
		if err != nil {
			return Token{}, err
		}

		s.currentToken = token
		s.initialized = true
	}
	if s.currentToken.IsExpired(s.clock) {
		token, err := renewToken(s.currentToken, s.sc, s.clock.Now())
		if err != nil {
			return Token{}, err
		}

		s.currentToken = token
	}
	return s.currentToken, nil
}

func newToken(sc SoundcloudApi, now time.Time) (Token, error) {
	tokenData, err := sc.Auth()
	if err != nil {
		return Token{}, err
	}

	t, err := NewTokenFromJsonData(tokenData, now)
	if err != nil {
		return Token{}, err
	}

	return t, nil
}

func renewToken(t Token, sc SoundcloudApi, now time.Time) (Token, error) {
	renewData, err := sc.Renew(t)
	if err != nil {
		return Token{}, err
	}

	newTok, err := NewTokenFromJsonData(renewData, now)
	if err != nil {
		return Token{}, err
	}

	return newTok, nil
}
