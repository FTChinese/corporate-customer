package controller

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "b2b.controller")

const claimsCtxKey = "claims"

type Doorkeeper struct {
	signingKey []byte
}

func NewDoorkeeper(key []byte) Doorkeeper {
	return Doorkeeper{
		signingKey: key,
	}
}

func (keeper Doorkeeper) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		ss, err := ParseBearer(authHeader)
		if err != nil {
			logger.Printf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := model.ParsePassportClaims(ss, keeper.signingKey)
		if err != nil {
			logger.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func (keeper Doorkeeper) CheckInviteeClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		ss, err := ParseBearer(authHeader)
		if err != nil {
			logger.Printf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := model.ParseInviteeClaims(ss, keeper.signingKey)
		if err != nil {
			logger.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func getPassportClaims(c echo.Context) model.PassportClaims {
	return c.Get(claimsCtxKey).(model.PassportClaims)
}

func getInviteeClaims(c echo.Context) model.InviteeClaims {
	return c.Get(claimsCtxKey).(model.InviteeClaims)
}
