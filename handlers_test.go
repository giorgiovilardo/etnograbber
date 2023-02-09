package main

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	expectedResponseBody := `{"ok":"okeione"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, HealthHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
	}
}
