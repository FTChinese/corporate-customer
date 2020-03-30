package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/views"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetForgotPassword show a form to collection user's email
func (router BarrierRouter) GetForgotPassword(c echo.Context) error {
	ctx := views.NewCtxBuilder().
		WithForm(views.NewResetLetterForm(admin.AccountForm{})).
		Build()
	return c.Render(http.StatusOK, "password_reset_email.html", ctx)
}

// PostForgotPassword handles sending email to help reset password.
func (router BarrierRouter) PostForgotPassword(c echo.Context) error {
	var af admin.AccountForm
	if err := c.Bind(&af); err != nil {
		return err
	}

	if ok := af.ValidatePasswordReset(); !ok {
		ctx := views.NewCtxBuilder().
			WithForm(views.NewResetLetterForm(af)).
			Build()

		return c.Render(http.StatusOK, "password_reset_email.html", ctx)
	}

	ctx := views.NewCtxBuilder().Set("done", true).Build()

	return c.Render(http.StatusOK, "password_reset_email.html", ctx)
}

func (router BarrierRouter) VerifyPasswordToken(c echo.Context) error {
	return nil
}

func (router BarrierRouter) GetResetPassword(c echo.Context) error {
	return nil
}

func (router BarrierRouter) PostResetPassword(c echo.Context) error {
	return nil
}
