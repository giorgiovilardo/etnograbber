package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	config := GetConfig()
	clock := RealClock{}
	httpSoundcloudApi := NewHttpSoundcloudApi(config)
	httpTokenRepository := NewHttpTokenRepository(clock, httpSoundcloudApi)
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))
	e.Use(loadTokenRepository(httpTokenRepository))
	e.Use(loadTrackRepository(httpSoundcloudApi))
	e.GET("/health", HealthHandler)
	e.GET("/:trackId", TrackHandler)
	e.GET("/:trackId/stream", StreamTrackHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

func loadTokenRepository(t TokenRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("tokenRepository", t)
			return next(c)
		}
	}
}

func loadTrackRepository(t TrackRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("trackRepository", t)
			return next(c)
		}
	}
}
