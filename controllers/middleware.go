package controllers

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = logrus.WithField("package", "controllers")

const sessionKey = "_ftc_b2b"
const loggedInKey = "account"

var SiteMap = struct {
	Home  string
	Login string
}{
	Home:  "/",
	Login: "/login",
}

func IgnoreFavicon(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/favicon.ico" {
			return nil
		}

		return next(c)
	}
}

func NotLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Infof("Path: %s", c.Path())

		sess, err := session.Get(sessionKey, c)

		if err != nil {
			return c.Redirect(http.StatusFound, SiteMap.Login)
		}

		logger.Infof("Session is new: %t", sess.IsNew)
		logger.Infof("Session values: %+v", sess.Values)

		_, ok := sess.Values[loggedInKey]
		// Not logged in
		if !ok {
			// Allow a non-logged-in user to access login page.
			if c.Path() == SiteMap.Login {
				return next(c)
			}
			// Redirect user to login from any other page.
			return c.Redirect(http.StatusFound, SiteMap.Login)
		}

		return next(c)
	}
}

func AlreadyLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(sessionKey, c)

		if err != nil {
			return c.Redirect(http.StatusFound, SiteMap.Login)
		}

		_, ok := sess.Values[loggedInKey]
		// Logged in
		if ok && c.Path() == SiteMap.Login {
			// Redirect user to home page if it is trying to access login page.
			return c.Redirect(http.StatusFound, SiteMap.Home)
		}

		return next(c)
	}
}
