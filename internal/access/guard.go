package access

import (
	"database/sql"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
)

type Guard struct {
	repo Env
}

func NewGuard(dbs db.ReadWriteMyDBs) Guard {
	return Guard{
		repo: NewEnv(dbs),
	}
}

func (g Guard) RequireToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := xhttp.GetAccessToken(c.Request())
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
