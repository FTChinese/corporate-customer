package controller

import (
	"errors"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/enum"
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

// WxSignUp handles a wechat user's intent to create a new email account
// and link it.
// Client should update login session since the reader.Passport is changed.
func (router ReaderRouter) WxSignUp(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxSignUp(
			claims.UnionID.String,
			c.Request().Body,
		)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

// WxLink lets a wechat logged-in user to link an existing email,
// or an email user to link to wechat.
// You need to get the linking target's account before sending request here.
// Client should update login session since the reader.Passport is changed.
func (router ReaderRouter) WxLink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxLink(
			claims.UnionID.String,
			c.Request().Body,
		)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	if resp.StatusCode != 204 {
		return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
	}

	fetchResp, err := router.apiClient.LoadAccountByUnionID(claims.UnionID.String)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, fetchResp)
}

// WxUnlink severs the link between wechat and email.
// Client should update login session since the reader.Passport is changed.
func (router ReaderRouter) WxUnlink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.apiClient.
		WxUnlink(
			claims.UnionID.String,
			c.Request().Body,
		)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	if resp.StatusCode != 204 {
		return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
	}

	var fResp fetch.Response
	switch claims.LoginMethod {
	case enum.LoginMethodEmail, enum.LoginMethodMobile:
		fResp, err = router.apiClient.LoadAccountByFtcID(claims.FtcID)
	case enum.LoginMethodWx:
		fResp, err = router.apiClient.LoadAccountByUnionID(claims.UnionID.String)
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, fResp)
}
