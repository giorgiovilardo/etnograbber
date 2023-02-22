package main

import (
	"errors"
)

type HttpTrackDataService struct {
	tr  TokenRepository
	tdr TrackDataRepository
}

func NewHttpTrackDataService(tr TokenRepository, tdr TrackDataRepository) *HttpTrackDataService {
	return &HttpTrackDataService{tr: tr, tdr: tdr}
}

func (t *HttpTrackDataService) GetTrackData(id int) (map[string]interface{}, error) {
	token, err := t.tr.GetToken()
	if err != nil {
		return nil, errors.New("token not available")
	}

	track, err := t.tdr.GetTrackData(token, id)
	if err != nil {
		return nil, errors.New("trackData not available")
	}

	return track, nil
}

type HttpCachedTrackService struct {
	c   TrackCache
	tr  TokenRepository
	trr TrackRepository
}

func NewHttpCachedTrackService(c TrackCache, tr TokenRepository, trr TrackRepository) *HttpCachedTrackService {
	return &HttpCachedTrackService{c: c, tr: tr, trr: trr}
}

func (t *HttpCachedTrackService) GetTrack(id int) ([]byte, error) {
	if t.c.Contains(id) {
		track, _ := t.c.Get(id)
		return track, nil
	}

	token, err := t.tr.GetToken()
	if err != nil {
		return nil, errors.New("token not available")
	}

	track, err := t.trr.GetTrack(token, id)
	if err != nil {
		return nil, errors.New("track not available")
	}

	t.c.Add(id, track)
	return track, nil
}
