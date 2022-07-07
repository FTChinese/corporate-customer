package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router StripeRouter) GetLatestInvoice(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.clients.
		Select(claims.Live).
		StripeSubsLatestInvoice(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) CouponOfLatestInvoice(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.clients.
		Select(claims.Live).
		CouponOfLatestSubsInvoice(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
