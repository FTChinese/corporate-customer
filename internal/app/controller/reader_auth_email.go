package controller

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
)

func (router ReaderRouter) EmailExists(c echo.Context) error {
	rawQuery := c.QueryString()

	resp, err := router.apiClient.EmailExists(rawQuery)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) handlePassport(c echo.Context, resp fetch.Response) error {
	// Forward error response back directly.
	if resp.StatusCode != 200 {
		return c.JSONBlob(resp.StatusCode, resp.Body)
	}

	// Use the bytes to create jwt token and then add the field to the bytes.
	var a reader.Account
	if err := json.Unmarshal(resp.Body, &a); err != nil {
		return render.NewInternalError(err.Error())
	}

	pp, err := router.guard.CreatePassport(a)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSON(resp.StatusCode, pp)
}

// EmailLogin handles login by email + password.
// Headers required to forward:
// * X-Client-Type: always `web`
// * X-Client-Version: user frontend or backend version?
// * X-User-Ip: retrieved from Context.RealIp
// * X-User-Agent: retrieved from Context.Request[User-Agent]
func (router ReaderRouter) EmailLogin(c echo.Context) error {

	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.apiClient.
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
func (router ReaderRouter) EmailSignUp(c echo.Context) error {

	defer c.Request().Body.Close()

	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.apiClient.EmailSignUp(c.Request().Body, header)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

func (router ReaderRouter) VerifyEmail(c echo.Context) error {
	token := c.Param("token")

	resp, err := router.apiClient.VerifyEmail(token)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
