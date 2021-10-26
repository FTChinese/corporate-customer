package controller

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/labstack/echo/v4"
)

const claimsCtxKey = "claims"

func getAdminClaims(c echo.Context) admin.PassportClaims {
	return c.Get(claimsCtxKey).(admin.PassportClaims)
}

func getReaderClaims(c echo.Context) reader.PassportClaims {
	return c.Get(claimsCtxKey).(reader.PassportClaims)
}
