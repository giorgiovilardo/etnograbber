package main

type SoundcloudApi interface {
	Auth() ([]byte, error)
	Renew(t Token) ([]byte, error)
}

type TokenRepository interface {
	GetToken() (Token, error)
}

type TrackCache interface {
	Add(key int, value []byte) (evicted bool)
	Contains(key int) bool
	Get(key int) (value []byte, ok bool)
}

type TrackRepository interface {
	GetTrack(t Token, id int) ([]byte, error)
}

type TrackDataRepository interface {
	GetTrackData(t Token, id int) (map[string]interface{}, error)
}

type TrackDataService interface {
	GetTrackData(id int) (map[string]interface{}, error)
}

type TrackService interface {
	GetTrack(id int) ([]byte, error)
}
