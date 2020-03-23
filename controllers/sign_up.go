package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/views"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router BarrierRouter) GetSignUp(c echo.Context) error {
	ctx := views.NewCtxBuilder().
		WithForm(views.NewSignUpForm(admin.SignUp{})).
		Build()

	return c.Render(http.StatusOK, "signup.html", ctx)
}

func (router BarrierRouter) PostSignUp(c echo.Context) error {
	var s admin.SignUp
	if err := c.Bind(&s); err != nil {
		return err
	}

	if ok := s.Sanitize().Validate(); !ok {
		ctx := views.NewCtxBuilder().
			WithForm(views.NewSignUpForm(s)).
			Build()

		return c.Render(http.StatusOK, "signup.html", ctx)
	}

	// TODO: Save signup data to db; retrieve the account.

	sess := createSession(c)
	sess.Values[loggedInKey] = s.Email
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Home)
}
