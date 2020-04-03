package controllers

import (
	"github.com/FTChinese/b2b/models/form"
	"github.com/FTChinese/b2b/views"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router BarrierRouter) GetSignUp(c echo.Context) error {
	ctx := views.NewCtxBuilder().
		WithForm(views.NewSignUpForm(form.AccountForm{})).
		Build()

	return c.Render(http.StatusOK, "signup.html", ctx)
}

func (router BarrierRouter) PostSignUp(c echo.Context) error {
	var af form.AccountForm
	if err := c.Bind(&af); err != nil {
		return err
	}

	if ok := af.ValidateSignUp(); !ok {
		ctx := views.NewCtxBuilder().
			WithForm(views.NewSignUpForm(af)).
			Build()

		return c.Render(http.StatusOK, "signup.html", ctx)
	}

	// TODO: Save signup data to db; retrieve the account.

	sess := createSession(c)
	sess.Values[loggedInKey] = af.Email
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Home)
}
