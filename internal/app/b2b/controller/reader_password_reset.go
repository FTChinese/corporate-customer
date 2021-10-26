package controller

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

// ResetPassword allow user to change password.
//
//	POST /users/password-reset
//
// Input:
// * token: string;
// * password: string.
func (router ReaderRouter) ResetPassword(c echo.Context) error {
	resp, err := router.apiClient.ResetPassword(c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// RequestPwResetLetter checks user's email and send a password
// reset letter if it is valid.
//
//	POST /users/password-reset/letter
//
// Input:
// * email: string;
// * sourceUrl: string.
//
// The footprint.Client headers are required.
func (router ReaderRouter) RequestPwResetLetter(c echo.Context) error {
	resp, err := router.apiClient.RequestPasswordResetLetter(c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// VerifyPwResetToken verifies a password reset link.
//
// 	GET /auth/password-reset/tokens/{token}
func (router ReaderRouter) VerifyPwResetToken(c echo.Context) error {
	token := c.Param("token")

	resp, err := router.apiClient.VerifyResetToken(token)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
