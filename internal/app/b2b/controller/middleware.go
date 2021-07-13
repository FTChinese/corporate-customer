package controller

import (
	admin2 "github.com/FTChinese/ftacademy/internal/pkg/admin"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http/httputil"
)

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
			log.Errorf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := admin2.ParsePassportClaims(ss, keeper.signingKey)
		if err != nil {
			log.Errorf("Error parsing JWT %v", err)
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
			log.Printf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := model2.ParseInviteeClaims(ss, keeper.signingKey)
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func getPassportClaims(c echo.Context) admin2.PassportClaims {
	return c.Get(claimsCtxKey).(admin2.PassportClaims)
}

func getInviteeClaims(c echo.Context) model2.InviteeClaims {
	return c.Get(claimsCtxKey).(model2.InviteeClaims)
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
