package controller

import (
	"github.com/FTChinese/b2b/internal/app/b2b/repository/products"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductRouter struct {
	repo products.Env
}

func NewProductRouter(repo products.Env) ProductRouter {
	return ProductRouter{repo: repo}
}

// ListProducts loads all products.
func (router ProductRouter) ListProducts(c echo.Context) error {
	prod, err := router.repo.LoadProducts()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, prod)
}
