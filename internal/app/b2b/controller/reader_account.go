package controller

import (
	"errors"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router ReaderRouter) LoadAccount(c echo.Context) error {
	claims := getReaderClaims(c)

	var resp fetch.Response
	var err error

	if claims.FtcID != "" {
		resp, err = router.apiClient.LoadAccountByFtcID(claims.FtcID)
	} else if claims.UnionID.Valid {
		resp, err = router.apiClient.LoadAccountByUnionID(claims.UnionID.String)
	} else {
		err = errors.New("unauthorized access")
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSONBlob(resp.StatusCode, resp.Body)
}

func (router ReaderRouter) LoadAccountWithJWT(c echo.Context) error {
	claims := getReaderClaims(c)

	var resp fetch.Response
	var err error

	if claims.FtcID != "" {
		resp, err = router.apiClient.LoadAccountByFtcID(claims.FtcID)
	} else if claims.UnionID.Valid {
		resp, err = router.apiClient.LoadAccountByUnionID(claims.UnionID.String)
	} else {
		err = errors.New("unauthorized access")
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

func (router ReaderRouter) UpdateEmail(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClient.
		UpdateEmail(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) RequestVerification(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.RequestEmailVerification(claims.FtcID)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) UpdateName(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		UpdateName(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) UpdatePassword(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		UpdatePassword(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) RequestMobileUpdateSMS(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		RequestMobileUpdateSMS(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) UpdateMobile(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		UpdateMobile(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) LoadAddress(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.
		LoadAddress(claims.FtcID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) UpdateAddress(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		UpdateAddress(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) LoadProfile(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.apiClient.
		LoadProfile(claims.FtcID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) UpdateProfile(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		UpdateProfile(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) WxSignUp(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxSignUp(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) WxLink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxLink(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) WxUnlink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxUnlink(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
