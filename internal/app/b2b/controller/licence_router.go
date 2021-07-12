package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LicenceRouter struct {
	repo subs.Env
}

func NewLicenceRouter(env subs.Env) LicenceRouter {
	return LicenceRouter{
		repo: env,
	}
}

func (router LicenceRouter) ListLicence(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	countCh, listCh := router.repo.AsyncCountLicences(claims.TeamID.String), router.repo.AsyncListExpLicence(claims.TeamID.String, page)

	countResult, listResult := <-countCh, <-listCh
	if listResult.Err != nil {
		return render.NewDBError(listResult.Err)
	}

	listResult.Total = countResult.Total
	return c.JSON(http.StatusOK, listResult)
}

func (router LicenceRouter) UpdateLicence(c echo.Context) error {
	return nil
}

// RevokeLicence unlinks a reader from a licence.
func (router LicenceRouter) RevokeLicence(c echo.Context) error {
	id := c.Param("id")

	claims := getPassportClaims(c)

	err := router.repo.RevokeLicence(id, claims.TeamID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
