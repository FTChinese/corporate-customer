package controllers

import (
	"github.com/FTChinese/corporate-customer/models"
	"github.com/FTChinese/corporate-customer/views"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SignInRouter struct {
	db   *sqlx.DB
	post postoffice.Postman
}

func NewSignInRouter(db *sqlx.DB, p postoffice.Postman) SignInRouter {
	return SignInRouter{db: db}
}

func (router SignInRouter) GetLogin(c echo.Context) error {
	data := views.BuildLoginPage(models.Login{})

	return c.Render(http.StatusOK, "login.html", data)
}
