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
	httpTrackDataService := NewHttpTrackDataService(httpTokenRepository, httpSoundcloudApi)
	trackCache, _ := lru.New[int, []byte](config.CacheSize)
	httpCachedTrackService := NewHttpCachedTrackService(trackCache, httpTokenRepository, httpSoundcloudApi)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: config.AllowedOrigins, AllowMethods: []string{http.MethodGet}}))
	e.GET("/health", HealthHandler)
	e.GET("/:trackId", TrackDataHandler(httpTrackDataService))
	e.GET("/:trackId/stream", TrackHandler(httpCachedTrackService))
	e.Logger.Fatal(e.Start(config.Address))
}
