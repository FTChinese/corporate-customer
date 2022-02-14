package b2b

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoadTeam retrieves a b2b team by the id specified
// in url path parameters.
func (router CMSRouter) LoadTeam(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	teamID := c.Param("id")

	t, err := router.repo.LoadTeam(teamID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, t)
}

func (router CMSRouter) LoadingAdminProfile(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	adminID := c.Param("id")

	profile, err := router.repo.LoadB2BAdminProfile(adminID)
	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, profile)
}
