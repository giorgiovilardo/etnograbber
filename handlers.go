package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}

func TrackHandler(c echo.Context) error {
	trackId, err := strconv.ParseInt(c.Param("trackId"), 10, 64)
	if err != nil {
		return trackIdNotNumberError(c)
	}

	token, err := c.Get("tokenRepository").(TokenRepository).GetToken()
	if err != nil {
		return tokenNotAvailableError(c)
	}

	track, err := c.Get("trackRepository").(TrackRepository).GetTrackData(token, int(trackId))
	if err != nil {
		return apiError(c, trackDataNotAvailable)
	}

	return c.JSON(http.StatusOK, track)
}

func StreamTrackHandler(c echo.Context) error {
	trackId, err := strconv.ParseInt(c.Param("trackId"), 10, 64)
	if err != nil {
		return trackIdNotNumberError(c)
	}

	token, err := c.Get("tokenRepository").(TokenRepository).GetToken()
	if err != nil {
		return tokenNotAvailableError(c)
	}

	trackReader, err := c.Get("trackRepository").(TrackRepository).GetTrack(token, int(trackId))
	if err != nil {
		return apiError(c, trackNotAvailable)
	}

	return c.Stream(http.StatusOK, "audio/mpeg", trackReader)
}

func apiError(c echo.Context, message string) error {
	return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": message})
}

func tokenNotAvailableError(c echo.Context) error {
	return apiError(c, tokenNotAvailable)
}

func trackIdNotNumberError(c echo.Context) error {
	return apiError(c, trackIdNotANumber)
}

const tokenNotAvailable = "token not available"
const trackIdNotANumber = "trackId not a number"
const trackDataNotAvailable = "trackData not available"
const trackNotAvailable = "track not available"
