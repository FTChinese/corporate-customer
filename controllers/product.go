package controllers

import (
	"github.com/FTChinese/b2b/repository"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductRouter struct {
	repo repository.Env
}

func NewProductRouter(repo repository.Env) ProductRouter {
	return ProductRouter{repo: repo}
}

// ListProducts loads all products.
func (router ProductRouter) ListProducts(c echo.Context) error {
	products, err := router.repo.LoadProducts()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, products)
}
