package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	config := GetConfig()
	clock := NewRealClock()
	httpSoundcloudApi := NewHttpSoundcloudApi(config)
	httpTokenRepository := NewHttpTokenRepository(clock, httpSoundcloudApi)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet},
	}))
	e.Use(addService(httpTokenRepository, "tokenRepository"))
	e.Use(addService(httpSoundcloudApi, "trackRepository"))
	e.GET("/health", HealthHandler)
	e.GET("/:trackId", TrackHandler)
	e.GET("/:trackId/stream", StreamTrackHandler)
	e.Logger.Fatal(e.Start(":5000"))
}

func addService[T any](t T, key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(key, t)
			return next(c)
		}
	}
}
