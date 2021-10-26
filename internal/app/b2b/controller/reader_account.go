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

//func (router ReaderRouter) UpdateEmail(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) RequestVerification(c echo.Context) error {
//
//}

//func (router ReaderRouter) UpdateName(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) UpdatePassword(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) SMSToModifyMobile(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) UpdateMobile(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) LoadAddress(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) UpdateAddress(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) LoadProfile(c echo.Context) error {
//
//}
//
//func (router ReaderRouter) UpdateProfile(c echo.Context) error {
//
//}
