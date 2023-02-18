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
	token, err := c.Get("tokenRepository").(TokenRepository).GetToken()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "token not available",
		})
	}

	trackId, err := strconv.ParseInt(c.Param("trackId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "trackId not a number",
		})
	}

	track, err := c.Get("trackRepository").(TrackRepository).GetTrackData(token, int(trackId))
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "trackData not available",
		})
	}

	return c.JSON(http.StatusOK, track)
}

func StreamTrackHandler(c echo.Context) error {
	token, err := c.Get("tokenRepository").(TokenRepository).GetToken()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "token not available",
		})
	}

	trackId, err := strconv.ParseInt(c.Param("trackId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "trackId not a number",
		})
	}

	trackReader, err := c.Get("trackRepository").(TrackRepository).GetTrack(token, int(trackId))
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "trackData not available",
		})
	}

	return c.Stream(http.StatusOK, "audio/mpeg", trackReader)
}
