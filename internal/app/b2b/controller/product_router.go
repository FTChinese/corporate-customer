package controller

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductRouter struct {
	apiClient api.Client
	logger    *zap.Logger
}

func NewProductRouter(client api.Client, logger *zap.Logger) ProductRouter {
	return ProductRouter{
		apiClient: client,
		logger:    logger,
	}
}

func (router ProductRouter) Paywall(c echo.Context) error {
	resp, err := router.apiClient.Paywall()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) ListStripePrices(c echo.Context) error {
	resp, err := router.apiClient.ListStripePrices()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) StripePrice(c echo.Context) error {
	id := c.QueryParam("id")
	resp, err := router.apiClient.StripePrice(id)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
