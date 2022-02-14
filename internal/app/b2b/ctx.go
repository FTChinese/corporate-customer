package b2b

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

func getAdminClaims(c echo.Context) admin.PassportClaims {
	return c.Get(xhttp.KeyCtxClaims).(admin.PassportClaims)
}
