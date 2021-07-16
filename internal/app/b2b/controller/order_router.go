package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/subs"
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
		repo: subs.NewEnv(dbs, logger),
	}
}

// CreateOrders creates orders an org purchased.
// Client should specify which plans are being subscribed,
// and how many copies.
func (router OrderRouter) CreateOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	return c.JSON(http.StatusOK, claims)
}

func (router OrderRouter) ListOrders(c echo.Context) error {
	claims := getPassportClaims(c)

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.JSON(http.StatusOK, claims)
}
