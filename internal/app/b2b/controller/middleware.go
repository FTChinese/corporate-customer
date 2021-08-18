package controller

import (
	"database/sql"
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/access"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/oauth"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/http/httputil"
)

const claimsCtxKey = "claims"

type JWTGuard struct {
	signingKey []byte
}

func NewJWTGuard(key []byte) JWTGuard {
	return JWTGuard{
		signingKey: key,
	}
}

func (g JWTGuard) getPassportClaims(req *http.Request) (admin.PassportClaims, error) {
	ss, err := oauth.GetBearerAuth(req.Header)
	if err != nil {
		log.Printf("Error parsing Authorization header: %v", err)
		return admin.PassportClaims{}, err
	}

	claims, err := admin.ParsePassportClaims(ss, g.signingKey)
	if err != nil {
		log.Printf("Error parsing JWT %v", err)
		return admin.PassportClaims{}, err
	}

	return claims, nil
}

func (g JWTGuard) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims, err := g.getPassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func (g JWTGuard) RequireTeamSet(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims, err := g.getPassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		if claims.TeamID.IsZero() {
			log.Printf("Team is not set")
			return render.NewUnauthorized("Organization team is required")
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func getPassportClaims(c echo.Context) admin.PassportClaims {
	return c.Get(claimsCtxKey).(admin.PassportClaims)
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

type OAuthGuard struct {
	repo access.Env
}

func NewOAuthGuard(dbs db.ReadWriteMyDBs) OAuthGuard {
	return OAuthGuard{
		repo: access.NewEnv(dbs),
	}
}

func (g OAuthGuard) RequireToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := oauth.GetToken(c.Request())
		if err != nil {
			log.Printf("Token not found: %s", err)
			return render.NewForbidden("Invalid access token")
		}

		o, err := g.repo.Load(token)
		if err != nil {
			if err == sql.ErrNoRows {
				return render.NewForbidden("Invalid access token")
			}
			return render.NewDBError(err)
		}

		if o.Expired() || !o.Active {
			log.Printf("Token %s is either expired or not active", token)
			return render.NewForbidden("The access token is expired or no longer active")
		}

		return next(c)
	}
}
