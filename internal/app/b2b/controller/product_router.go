package controller

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/repository/api"
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
