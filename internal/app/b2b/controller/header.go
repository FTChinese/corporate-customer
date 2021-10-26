package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func collectClientHeader(c echo.Context, version string) http.Header {
	h := http.Header{}
	h.Set("X-Client-Type", "web")
	h.Set("X-Client-Version", version)
	h.Set("X-User-Ip", c.RealIP())
	h.Set("X-User-Agent", c.Request().UserAgent())

	return h
}
