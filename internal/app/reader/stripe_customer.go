package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router StripeRouter) CreateCustomer(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.StripeNewCustomer(claims)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetCustomer(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeGetSubs(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetCusDefaultPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	cusID := c.Param("id")

	resp, err := router.apiClient.StripeCustomerDefaultPaymentMethod(claims, cusID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
