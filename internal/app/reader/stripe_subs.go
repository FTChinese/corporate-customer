package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router StripeRouter) CreateSubs(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClient.StripeNewSubs(claims, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetSubs(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeGetSubs(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) UpdateSubs(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	defer c.Request().Body.Close()

	resp, err := router.apiClient.StripeUpdateSubsOf(claims, subsID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) RefreshSubs(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeRefreshSubsOf(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) CancelSubs(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeCancelSubsOf(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) ReactivateSubs(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeReactivateSubsOf(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetSubsDefaultPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	subsID := c.Param("id")

	resp, err := router.apiClient.StripeSubsDefaultPaymentMethod(claims, subsID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) ListPaymentMethods(c echo.Context) error {
	claims := getReaderClaims(c)
	rawQuery := c.QueryString()

	resp, err := router.apiClient.StripeListPaymentMethods(claims, rawQuery)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router StripeRouter) GetPaymentMethod(c echo.Context) error {
	claims := getReaderClaims(c)
	pmID := c.Param("id")

	resp, err := router.apiClient.StripePaymentMethodOf(claims, pmID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
