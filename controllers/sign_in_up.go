package controllers

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/views"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SignInRouter struct {
	db   *sqlx.DB
	post postoffice.Postman
}

func NewSignInRouter(db *sqlx.DB, p postoffice.Postman) SignInRouter {
	return SignInRouter{db: db, post: p}
}

func (router SignInRouter) GetLogin(c echo.Context) error {
	data := views.BuildLoginPage(admin.Login{})

	return c.Render(http.StatusOK, "login.html", data)
}

func (router SignInRouter) PostLogin(c echo.Context) error {
	var l admin.Login
	if err := c.Bind(&l); err != nil {
		return err
	}

	if ok := l.Sanitize().Validate(); !ok {
		data := views.BuildLoginPage(l)

		return c.Render(http.StatusOK, "login.html", data)
	}

	sess, _ := session.Get(sessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
	}

	sess.Values[loggedInKey] = l.Email
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Home)
}

func (router SignInRouter) LogOut(c echo.Context) error {
	sess, _ := session.Get(sessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1, // A zero or negative number will expire the cookie immediately. If both Expires and Max-Age are set, Max-Age has precedence.
		HttpOnly: false,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Login)
}

func (router SignInRouter) GetForgotPassword(c echo.Context) error {
	return nil
}

func (router SignInRouter) PostForgotPassword(c echo.Context) error {
	return nil
}

func (router SignInRouter) VerifyPasswordToken(c echo.Context) error {
	return nil
}

func (router SignInRouter) GetResetPassword(c echo.Context) error {
	return nil
}

func (router SignInRouter) PostResetPassword(c echo.Context) error {
	return nil
}
