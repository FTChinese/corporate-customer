package controllers

import (
	"github.com/FTChinese/b2b/repository"
	"github.com/labstack/echo/v4"
)

type LicenceRouter struct {
	repo repository.Env
}

func NewLicenceRouter(env repository.Env) LicenceRouter {
	return LicenceRouter{
		repo: env,
	}
}

func (router LicenceRouter) ListLicence(c echo.Context) error {
	return nil
}

func (router LicenceRouter) UpdateLicence(c echo.Context) error {
	return nil
}

func (router LicenceRouter) RevokeLicence(c echo.Context) error {
	return nil
}
