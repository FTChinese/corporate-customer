package controllers

import (
	"github.com/FTChinese/b2b/models/form"
	"github.com/FTChinese/b2b/views"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BarrierRouter struct {
	db   *sqlx.DB
	post postoffice.PostOffice
}

func NewBarrierRouter(db *sqlx.DB, p postoffice.PostOffice) BarrierRouter {
	return BarrierRouter{db: db, post: p}
}

func (router BarrierRouter) GetLogin(c echo.Context) error {
	ctx := views.NewCtxBuilder().
		WithForm(views.NewLoginForm(form.AccountForm{})).
		Build()

	return c.Render(http.StatusOK, "login.html", ctx)
}

func createSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get(sessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
	}

	return sess
}

func (router BarrierRouter) PostLogin(c echo.Context) error {
	var af form.AccountForm
	if err := c.Bind(&af); err != nil {
		return err
	}

	if ok := af.ValidateLogin(); !ok {
		ctx := views.NewCtxBuilder().
			WithForm(views.NewLoginForm(af)).
			Build()

		return c.Render(http.StatusOK, "login.html", ctx)
	}

	sess := createSession(c)
	sess.Values[loggedInKey] = af.Email
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Home)
}

func (router BarrierRouter) LogOut(c echo.Context) error {
	sess, _ := session.Get(sessionKey, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1, // A zero or negative number will expire the cookie immediately. If both Expires and Max-Age are set, Max-Age has precedence.
		HttpOnly: false,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, SiteMap.Login)
}
