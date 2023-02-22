package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}

func TrackDataHandler(s TrackDataService) func(c echo.Context) error {
	return func(c echo.Context) error {
		trackId, err := strconv.Atoi(c.Param("trackId"))
		if err != nil {
			return apiError(c, trackIdNotANumber)
		}

		track, err := s.GetTrackData(trackId)
		if err != nil {
			return apiError(c, err.Error())
		}

		return c.JSON(http.StatusOK, track)
	}
}

func TrackHandler(s TrackService) func(c echo.Context) error {
	return func(c echo.Context) error {
		trackId, err := strconv.Atoi(c.Param("trackId"))
		if err != nil {
			return apiError(c, trackIdNotANumber)
		}

		track, err := s.GetTrack(trackId)
		if err != nil {
			return apiError(c, err.Error())
		}

		return c.Stream(http.StatusOK, "audio/mpeg", bytes.NewReader(track))
	}
}

func apiError(c echo.Context, message string) error {
	return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": message})
}

const trackIdNotANumber = "trackId not a number"
