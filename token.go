package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token" validate:"required"`
	ExpiresIn    int    `json:"expires_in" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
	Scope        string `json:"scope" validate:"eq="`
	TokenType    string `json:"token_type" validate:"required,eq=bearer"`
	ExpiresAt    time.Time
}

func (t Token) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func NewTokenFromJsonData(tokenData []byte, now time.Time) (Token, error) {
	var t Token
	if err := json.Unmarshal(tokenData, &t); err != nil {
		return Token{}, err
	}
	t.ExpiresAt = now.Add(time.Second * time.Duration(t.ExpiresIn))
	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return Token{}, err
	}
	return t, nil
}
