package controller

import (
	api2 "github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type ReaderRouter struct {
	guard     reader.JWTGuard
	apiClient api2.Client
	wxApp     config.WechatApp
	version   string
}

func NewReaderRouter(client api2.Client, version string) ReaderRouter {
	return ReaderRouter{
		guard: reader.NewJWTGuard(
			config.
				MustGetReaderAppKey().
				GetJWTKey(),
		),
		apiClient: client,
		wxApp:     config.MustWxWebApp(),
		version:   version,
	}
}

func (router ReaderRouter) collectClientHeader(c echo.Context) http.Header {
	return api2.NewHeaderBuilder().
		WithPlatformWeb().
		WithClientVersion(router.version).
		WithUserIP(c.RealIP()).
		WithUserAgent(c.Request().UserAgent()).
		Build()
}

func (router ReaderRouter) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims, err := router.guard.RetrievePassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}
