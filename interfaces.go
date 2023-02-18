package main

import (
	"io"
	"time"
)

type Clock interface {
	Now() time.Time
}

type SoundcloudApi interface {
	Auth() ([]byte, error)
	Renew(t Token) ([]byte, error)
}

type TokenRepository interface {
	GetToken() (Token, error)
}

type TrackRepository interface {
	GetTrackData(t Token, id int) (map[string]interface{}, error)
	GetTrack(t Token, id int) (io.Reader, error)
}
