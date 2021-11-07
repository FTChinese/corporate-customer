package controller

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/api"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router ReaderRouter) WxRequestCode(c echo.Context) error {
	sess, err := router.apiClient.WxOAuthSession(router.wxApp.AppID)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSON(http.StatusOK, sess)
}

func (router ReaderRouter) WxLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	header := router.collectClientHeader(c)

	resp, err := router.apiClient.WxLogin(
		router.wxApp.AppID,
		c.Request().Body,
		header)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	if resp.StatusCode != 200 {
		return c.JSONBlob(resp.StatusCode, resp.Body)
	}

	var sess api.WxLoginSession
	if err := json.Unmarshal(resp.Body, &sess); err != nil {
		return render.NewInternalError(err.Error())
	}

	// Now that we have wechat union id, use the id to fetch account data.
	resp, err = router.apiClient.LoadAccountByUnionID(sess.UnionID)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

func (router ReaderRouter) WxRefresh(c echo.Context) error {
	defer c.Request().Body.Close()

	header := router.collectClientHeader(c)

	resp, err := router.apiClient.WxRefresh(
		router.wxApp.AppID,
		c.Request().Body,
		header)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
