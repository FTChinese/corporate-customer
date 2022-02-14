package b2b

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/repository/adminrepo"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
)

// AdminRouter defines what an admin can do before login.
type AdminRouter struct {
	guard  admin.JWTGuard
	repo   adminrepo.Env
	post   postman.Postman
	logger *zap.Logger
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(dbs db.ReadWriteMyDBs, p postman.Postman, logger *zap.Logger) AdminRouter {
	return AdminRouter{
		guard: admin.NewJWTGuard(
			config.
				MustGetB2BAppKey().
				GetJWTKey(),
		),
		repo:   adminrepo.NewEnv(dbs, logger),
		post:   p,
		logger: logger,
	}
}

func (router AdminRouter) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
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

func (router AdminRouter) RequireTeamSet(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims, err := router.guard.RetrievePassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		if claims.TeamID.IsZero() {
			log.Printf("Team is not set")
			return render.NewUnauthorized("Organization team is required")
		}

		c.Set(xhttp.KeyCtxClaims, claims)
		return next(c)
	}
}
