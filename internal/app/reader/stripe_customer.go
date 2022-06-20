package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router StripeRouter) CreateCustomer(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.clients.
		Select(claims.Live).
		StripeNewCustomer(claims)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetCustomer(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.clients.
		Select(claims.Live).
		StripeGetCustomer(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetCusDefaultPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	cusID := c.Param("id")

	resp, err := router.clients.
		Select(claims.Live).
		StripeCusDefaultPaymentMethod(claims, cusID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) SetCusDefaultPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	cusID := c.Param("id")

	defer c.Request().Body.Close()

	resp, err := router.clients.
		Select(claims.Live).
		StripeSetCusDefaultPaymentMethod(claims, cusID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ListCusPaymentMethods list a customer's payment methods with pagination.
func (router StripeRouter) ListCusPaymentMethods(c echo.Context) error {
	claims := getReaderClaims(c)
	cusID := c.Param("id")

	resp, err := router.clients.
		Select(claims.Live).
		StripeListCusPaymentMethods(claims, cusID, c.QueryParams())

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
