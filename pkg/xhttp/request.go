package xhttp

import (
	"github.com/FTChinese/ftacademy/pkg/conv"
	"github.com/labstack/echo/v4"
)

func GetQueryLive(c echo.Context) bool {
	return conv.DefaultTrue(c.QueryParam("live"))
}

func GetQueryRefresh(c echo.Context) bool {
	return conv.DefaultFalse(c.QueryParam("refresh"))
}
