package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/checkout"
	"github.com/FTChinese/ftacademy/pkg/db"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type OrderRouter struct {
	repo   subs.Env
	logger *zap.Logger
}

func NewOrderRouter(dbs db.ReadWriteMyDBs, logger *zap.Logger) OrderRouter {
	return OrderRouter{
		repo:   subs.NewEnv(dbs, logger),
		logger: logger,
	}
}

// CreateOrders creates orders an org purchased.
// input: input.ShoppingCart.
func (router OrderRouter) CreateOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	var cart checkout.ShoppingCart
	if err := c.Bind(&cart); err != nil {
		return render.NewBadRequest(err.Error())
	}

	bo, err := router.repo.CreateOrder(cart, claims)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, bo)
}

func (router OrderRouter) ListOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	list, err := router.repo.ListOrders(claims.TeamID.String, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, list)
}

func (router OrderRouter) LoadOrder(c echo.Context) error {
	claims := getPassportClaims(c)

	id := c.QueryParam("id")

	o, err := router.repo.LoadOrderDetails(admin.AccessRight{
		RowID:  id,
		TeamID: claims.TeamID.String,
	})

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, o)
}
