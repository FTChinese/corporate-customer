package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

// GetPaymentMethod loads or refreshes a payment method.
// If query refresh=true presents, the API will fetch directly from Stripe.
func (router StripeRouter) GetPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	pmID := c.Param("id")

	resp, err := router.apiClient.StripePaymentMethodOf(
		claims,
		pmID,
		c.QueryParams())

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// CreateSetupIntent create a new setup intent so that client could
// call Stripe.js's payment element to create a new payment method
// for future usage.
func (router StripeRouter) CreateSetupIntent(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClient.StripeCreateSetupIntent(
		claims,
		c.Request().Body,
	)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// GetSetupIntent loads a setup intent.
func (router StripeRouter) GetSetupIntent(c echo.Context) error {
	claims := getReaderClaims(c)
	id := c.Param("id")

	resp, err := router.apiClient.StripeGetSetupIntent(
		claims,
		id,
		c.QueryParams(),
	)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// GetSetupPaymentMethod loads a payment method under a setup intent.
func (router StripeRouter) GetSetupPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	id := c.Param("id")

	resp, err := router.apiClient.StripeGetSetupIntent(
		claims,
		id,
		c.QueryParams(),
	)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
