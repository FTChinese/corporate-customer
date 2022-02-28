package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router Router) CreateWxOrder(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClient.WxPayDesktop(claims, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) CreateAliOrder(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClient.AliPayDesktop(
		claims,
		c.Request().Body,
	)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body,
	)
}

func (router Router) VerifyFtcOrder(c echo.Context) error {
	claims := getReaderClaims(c)

	orderID := c.Param("id")

	resp, err := router.apiClient.VerifyPaymentResult(
		claims,
		orderID,
	)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body,
	)
}
