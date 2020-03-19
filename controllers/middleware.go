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

func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Infof("Path: %s", c.Path())

		if c.Path() == "/favicon.ico" {
			return nil
		}

		sess, err := session.Get(sessionKey, c)

		if err != nil {
			return c.Redirect(http.StatusFound, SiteMap.Login)
		}

		logger.Infof("Session is new: %t", sess.IsNew)
		logger.Infof("Session values: %+v", sess.Values)

		_, ok := sess.Values[loggedInKey]
		if !ok {
			return c.Redirect(http.StatusFound, SiteMap.Login)
		}

		return next(c)
	}
}
