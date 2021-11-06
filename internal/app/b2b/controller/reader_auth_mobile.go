package controller

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
)

// RequestMobileLoginSMS sends an SMS to user for login.
// Input:
// * mobile: string
func (router ReaderRouter) RequestMobileLoginSMS(c echo.Context) error {
	resp, err := router.apiClient.RequestLoginSMS(c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// VerifyMobileLoginSMS verifies a code sent to user mobile devices.
// * mobile: string - the mobile number used for login
// * code: string - the SMS cod of this session
// Returns {id: string | null} containing uuid.
// If
func (router ReaderRouter) VerifyMobileLoginSMS(c echo.Context) error {
	resp, err := router.apiClient.VerifyLoginSMS(c.Request().Body)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	if resp.StatusCode != 200 {
		return c.JSONBlob(resp.StatusCode, resp.Body)
	}

	log.Printf("SMS login verification result %s", resp.Body)

	var found reader.SearchResult
	if err := json.Unmarshal(resp.Body, &found); err != nil {
		return render.NewInternalError(err.Error())
	}

	if found.ID.IsZero() {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "Mobile is not signed up",
			Field:   "mobile",
			Code:    render.CodeMissing,
		})
	}

	resp, err = router.apiClient.LoadAccountByFtcID(found.ID.String)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

// MobileLinkExistingEmail authenticates an existing email account,
// and link to the provided mobile phone.
//
// Input:
// * email: string;
// * password: string;
// * mobile: string;
func (router ReaderRouter) MobileLinkExistingEmail(c echo.Context) error {
	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.apiClient.MobileLinkExistingEmail(c.Request().Body, header)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}

// MobileSignUp creates a new mobile account.
// Input:
// * mobile: string;
func (router ReaderRouter) MobileSignUp(c echo.Context) error {

	header := router.collectClientHeader(c)

	log.Printf("Real IP: %v", header)

	resp, err := router.apiClient.MobileSignUp(c.Request().Body, header)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return router.handlePassport(c, resp)
}
