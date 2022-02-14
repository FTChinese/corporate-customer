package reader

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Router struct {
	guard     reader.JWTGuard
	apiClient api.Client
	wxApp     config.WechatApp
	version   string
}

func NewReaderRouter(client api.Client, version string) Router {
	return Router{
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

func (router Router) collectClientHeader(c echo.Context) http.Header {
	return api.NewHeaderBuilder().
		WithPlatformWeb().
		WithClientVersion(router.version).
		WithUserIP(c.RealIP()).
		WithUserAgent(c.Request().UserAgent()).
		Build()
}

func (router Router) handlePassport(c echo.Context, resp fetch.Response) error {
	// Forward error response back directly.
	if resp.StatusCode != 200 {
		return c.JSONBlob(resp.StatusCode, resp.Body)
	}

	// Use the bytes to create jwt token and then add the field to the bytes.
	var a reader.Account
	if err := json.Unmarshal(resp.Body, &a); err != nil {
		return render.NewInternalError(err.Error())
	}

	pp, err := router.guard.CreatePassport(a)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSON(resp.StatusCode, pp)
}

func getReaderClaims(c echo.Context) reader.PassportClaims {
	return c.Get(xhttp.KeyCtxClaims).(reader.PassportClaims)
}

func (router Router) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims, err := router.guard.RetrievePassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(xhttp.KeyCtxClaims, claims)
		return next(c)
	}
}
