package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToken_IsExpired(t *testing.T) {
	pastToken := Token{ExpiresAt: time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC)}
	futureToken := Token{ExpiresAt: time.Date(2921, 8, 25, 8, 30, 0, 0, time.UTC)}
	assert.True(t, pastToken.IsExpired(NewRealClock()))
	assert.False(t, futureToken.IsExpired(NewRealClock()))
}

func TestNewTokenFromJsonData(t *testing.T) {
	tests := []struct {
		name       string
		jsonRepr   []byte
		time       time.Time
		want       Token
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "should parse with expected json",
			jsonRepr: []byte(`{"access_token":"miao","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`),
			time:     time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC),
			want: Token{
				AccessToken:  "miao",
				ExpiresIn:    3599,
				ExpiresAt:    time.Date(2021, 8, 25, 9, 29, 59, 0, time.UTC),
				RefreshToken: "bau",
				Scope:        "",
				TokenType:    "bearer",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "should fail with wrong field types",
			jsonRepr:   []byte(`{"access_token":33,"expires_in":"z","refresh_token":"bau","scope":"","token_type":"bearer"}`),
			time:       time.Date(2021, 8, 25, 8, 0, 0, 0, time.UTC),
			want:       Token{},
			wantErr:    true,
			wantErrMsg: "json: cannot unmarshal number into Go struct field Token.access_token of type string",
		},
		{
			name:       "should fail with missing field",
			jsonRepr:   []byte(`{"expires_in":3599,"refresh_token":"bau","scope":"","token_type":"bearer"}`),
			time:       time.Date(2021, 8, 25, 8, 0, 0, 0, time.UTC),
			want:       Token{},
			wantErr:    true,
			wantErrMsg: "Key: 'Token.AccessToken' Error:Field validation for 'AccessToken' failed on the 'required' tag",
		},
		{
			name:       "should fail with different type than bearer",
			jsonRepr:   []byte(`{"access_token":"miao","expires_in":3599,"refresh_token":"bau","scope":"","token_type":"orso"}`),
			time:       time.Date(2021, 8, 25, 8, 0, 0, 0, time.UTC),
			want:       Token{},
			wantErr:    true,
			wantErrMsg: "Key: 'Token.TokenType' Error:Field validation for 'TokenType' failed on the 'eq' tag",
		},
		{
			name:       "should fail if scope is not empty",
			jsonRepr:   []byte(`{"access_token":"miao","expires_in":3599,"refresh_token":"bau","scope":"z","token_type":"bearer"}`),
			time:       time.Date(2021, 8, 25, 8, 0, 0, 0, time.UTC),
			want:       Token{},
			wantErr:    true,
			wantErrMsg: "Key: 'Token.Scope' Error:Field validation for 'Scope' failed on the 'eq' tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTokenFromJsonData(tt.jsonRepr, tt.time)
			if err != nil && tt.wantErr {
				assert.Equalf(t, tt.wantErrMsg, err.Error(), "tokenFromJson() error = %v, wantErrMsg %v", err, tt.wantErrMsg)
				return
			}
			assert.Equalf(t, tt.want, got, "NewTokenFromJsonData(%v, %v)", tt.jsonRepr, tt.time)
		})
	}
}
