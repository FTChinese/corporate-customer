package controllers

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReadersRouter struct {
	db   *sqlx.DB
	post postoffice.Postman
}

func NewReadersRouter(db *sqlx.DB, p postoffice.Postman) ReadersRouter {
	return ReadersRouter{
		db:   db,
		post: p,
	}
}

func (router ReadersRouter) GetUserList(c echo.Context) error {
	return c.Render(http.StatusOK, "reader-list.html", nil)
}

func (router ReadersRouter) GetReaderProfile(c echo.Context) error {
	return c.Render(http.StatusOK, "reader-profile.html", nil)
}
