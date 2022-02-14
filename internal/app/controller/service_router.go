package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
	"net/http"
)

func GenerateQRImage(c echo.Context) error {
	v := c.QueryParam("url")

	var png []byte
	png, err := qrcode.Encode(v, qrcode.Medium, 256)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Blob(http.StatusOK, "image/png", png)
}
