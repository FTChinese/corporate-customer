package controllers

import (
	"errors"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

var logger = logrus.WithField("package", "controllers")

const sessionKey = "_ftc_b2b"
const loggedInKey = "account"

func IgnoreFavicon(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/favicon.ico" {
			return nil
		}

		return next(c)
	}
}

// RequireLoggedIn router prevents access is user is not logged in.
func RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
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

// RedirectIfLoggedIn hides pages that are only available if user is not logged in, such as login, sign-up, forgot password, etc..
func RedirectIfLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
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

// ParseBearer extracts Authorization header.
// Authorization: Bearer 19c7d9016b68221cc60f00afca7c498c36c361e3
func ParseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	s := strings.SplitN(authHeader, " ", 2)

	bearerExists := (len(s) == 2) && (strings.ToLower(s[0]) == "bearer")

	if !bearerExists {
		return "", errors.New("bearer not found")
	}

	return s[1], nil
}

func CheckJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		ss, err := ParseBearer(authHeader)
		if err != nil {
			log.Printf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := admin.ParseJWT(ss)
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set("claims", claims)
		return next(c)
	}
}

func getAccountClaims(c echo.Context) admin.AccountClaims {
	return c.Get("claims").(admin.AccountClaims)
}

func DumpRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dump, err := httputil.DumpRequest(c.Request(), false)
		if err != nil {
			log.Print(err)
		}

		log.Printf(string(dump))

		return next(c)
	}
}
