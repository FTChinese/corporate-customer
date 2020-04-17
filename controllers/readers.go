package controllers

import (
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/labstack/echo/v4"
)

type ReadersRouter struct {
	repo repository.Env
	post postoffice.PostOffice
}

func NewReaderRouter(env repository.Env, post postoffice.PostOffice) ReadersRouter {
	return ReadersRouter{
		repo: env,
		post: post,
	}
}

func NewReadersRouter(env repository.Env, p postoffice.PostOffice) ReadersRouter {
	return ReadersRouter{
		repo: env,
		post: p,
	}
}

func (router ReadersRouter) SignUp(c echo.Context) error {
	return nil
}

func (router ReadersRouter) VerifyInvitation(c echo.Context) error {
	return nil
}

func (router ReadersRouter) Accept(c echo.Context) error {
	return nil
}
