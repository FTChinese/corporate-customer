package reader

import (
	"errors"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

func (router Router) LoadAccount(c echo.Context) error {
	claims := getReaderClaims(c)

	var resp fetch.Response
	var err error

	if claims.FtcID != "" {
		resp, err = router.clients.Select(true).LoadAccountByFtcID(claims.FtcID)
	} else if claims.UnionID.Valid {
		resp, err = router.clients.Select(true).LoadAccountByUnionID(claims.UnionID.String)
	} else {
		err = errors.New("unauthorized access")
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSONBlob(resp.StatusCode, resp.Body)
}

func (router Router) LoadAccountWithJWT(c echo.Context) error {
	claims := getReaderClaims(c)

	var resp fetch.Response
	var err error

	if claims.FtcID != "" {
		resp, err = router.clients.Select(true).LoadAccountByFtcID(claims.FtcID)
	} else if claims.UnionID.Valid {
		resp, err = router.clients.Select(true).LoadAccountByUnionID(claims.UnionID.String)
	} else {
		err = errors.New("unauthorized access")
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

func (router Router) UpdateEmail(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.clients.Select(true).
		UpdateEmail(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) RequestVerification(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.clients.Select(true).RequestEmailVerification(claims.FtcID)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) UpdateName(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		UpdateName(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) UpdatePassword(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		UpdatePassword(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) RequestMobileUpdateSMS(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		RequestMobileUpdateSMS(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) UpdateMobile(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		UpdateMobile(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) LoadAddress(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.clients.Select(true).
		LoadAddress(claims.FtcID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) UpdateAddress(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		UpdateAddress(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) LoadProfile(c echo.Context) error {
	claims := getReaderClaims(c)

	resp, err := router.clients.Select(true).
		LoadProfile(claims.FtcID)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router Router) UpdateProfile(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
		UpdateProfile(claims.FtcID, c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// WxSignUp handles a wechat user's intent to create a new email account
// and link it.
// Client should update login session since the reader.Passport is changed.
func (router Router) WxSignUp(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
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
func (router Router) WxLink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
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

	fetchResp, err := router.clients.Select(true).LoadAccountByUnionID(claims.UnionID.String)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, fetchResp)
}

// WxUnlink severs the link between wechat and email.
// Client should update login session since the reader.Passport is changed.
func (router Router) WxUnlink(c echo.Context) error {
	claims := getReaderClaims(c)

	defer c.Request().Body.Close()
	resp, err := router.clients.Select(true).
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
		fResp, err = router.clients.Select(true).LoadAccountByFtcID(claims.FtcID)
	case enum.LoginMethodWx:
		fResp, err = router.clients.Select(true).LoadAccountByUnionID(claims.UnionID.String)
	}

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, fResp)
}
