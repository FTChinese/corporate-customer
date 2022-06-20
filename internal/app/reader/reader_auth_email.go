package reader

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
)

func (router Router) EmailExists(c echo.Context) error {
	rawQuery := c.QueryString()

	resp, err := router.clients.Select(true).EmailExists(rawQuery)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// EmailLogin handles login by email + password.
// Headers required to forward:
// * X-Client-Type: always `web`
// * X-Client-Version: user frontend or backend version?
// * X-User-Ip: retrieved from Context.RealIp
// * X-User-Agent: retrieved from Context.Request[User-Agent]
func (router Router) EmailLogin(c echo.Context) error {

	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.clients.Select(true).
		EmailLogin(c.Request().Body, header)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

// EmailSignUp create a new account for a user.
//
// 	POST /users/signup
//
// Input:
// * email: string;
// * password: string;
// * mobile?: string;
// * sourceUrl: string From which site the request is sent. Not required for mobile apps.
//
// The footprint.Client headers are required.
func (router Router) EmailSignUp(c echo.Context) error {

	defer c.Request().Body.Close()

	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.clients.Select(true).EmailSignUp(c.Request().Body, header)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

func (router Router) VerifyEmail(c echo.Context) error {
	token := c.Param("token")

	resp, err := router.clients.Select(true).VerifyEmail(token)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
