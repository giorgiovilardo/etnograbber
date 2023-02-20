package main

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	config := GetConfig()
	clock := NewRealClock()
	httpSoundcloudApi := NewHttpSoundcloudApi(config)
	httpTokenRepository := NewHttpTokenRepository(clock, httpSoundcloudApi)
	trackCache, _ := lru.New[int, []byte](30)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet},
	}))
	e.GET("/health", HealthHandler)
	e.GET("/:trackId", TrackDataHandler(httpTokenRepository, httpSoundcloudApi))
	e.GET("/:trackId/stream", TrackHandler(trackCache, httpTokenRepository, httpSoundcloudApi))
	e.Logger.Fatal(e.Start(":5000"))
}
