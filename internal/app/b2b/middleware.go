package b2b

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http/httputil"
)

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

//const sessionKey = "_ftc_b2b"
//const loggedInKey = "account"

//func IgnoreFavicon(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		if c.Path() == "/favicon.ico" {
//			return nil
//		}
//
//		return next(c)
//	}
//}

// RequireLoggedIn router prevents access is user is not logged in.
//func RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		logger.Infof("Path: %s", c.Path())
//
//		sess, err := session.Get(sessionKey, c)
//
//		if err != nil {
//			return c.Redirect(http.StatusFound, SiteMap.Login)
//		}
//
//		logger.Infof("Session is new: %t", sess.IsNew)
//		logger.Infof("Session values: %+v", sess.Values)
//
//		_, ok := sess.Values[loggedInKey]
//		// Not logged in
//		if !ok {
//			// Allow a non-logged-in user to access login page.
//			if c.Path() == SiteMap.Login {
//				return next(c)
//			}
//			// Redirect user to login from any other page.
//			return c.Redirect(http.StatusFound, SiteMap.Login)
//		}
//
//		return next(c)
//	}
//}

// RedirectIfLoggedIn hides pages that are only available if user is not logged in, such as login, sign-up, forgot password, etc..
//func RedirectIfLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		sess, err := session.Get(sessionKey, c)
//
//		if err != nil {
//			return c.Redirect(http.StatusFound, SiteMap.Login)
//		}
//
//		_, ok := sess.Values[loggedInKey]
//		// Logged in
//		if ok && c.Path() == SiteMap.Login {
//			// Redirect user to home page if it is trying to access login page.
//			return c.Redirect(http.StatusFound, SiteMap.Home)
//		}
//
//		return next(c)
//	}
//}
