package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router Router) LoadMembership(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.Membership(claims)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) ClaimAddon(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.ClaimAddOn(claims)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
