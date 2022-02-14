package controller

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http/httputil"
)

func DumpRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dump, err := httputil.DumpRequest(c.Request(), false)
		if err != nil {
			log.Print(err)
		}

		log.Printf("\n------Dump request start------\n%s------Dump request end------", string(dump))

		return next(c)
	}
}

func NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := c.Response().Header()
		h.Add("Cache-Control", "no-cache")
		h.Add("Cache-Control", "no-store")
		h.Add("Cache-Control", "must-revalidate")
		h.Add("Pragma", "no-cache")
		return next(c)
	}
}
