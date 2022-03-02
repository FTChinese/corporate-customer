package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router Router) RefreshIAP(c echo.Context) error {
	claims := getReaderClaims(c)

	originalTxID := c.Param("id")

	resp, err := router.apiClient.RefreshIAP(
		claims,
		originalTxID,
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
