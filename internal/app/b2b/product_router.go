package b2b

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProductRouter struct {
	apiClient api.Client // Deprecated
	clients   api.Clients
	logger    *zap.Logger
}

func NewProductRouter(clients api.Clients, logger *zap.Logger) ProductRouter {
	return ProductRouter{
		clients: clients,
		logger:  logger,
	}
}

func (router ProductRouter) Paywall(c echo.Context) error {
	live := xhttp.GetQueryLive(c)

	resp, err := router.clients.Select(live).Paywall()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) StripePrice(c echo.Context) error {
	id := c.Param("id")
	live := xhttp.GetQueryLive(c)

	resp, err := router.clients.
		Select(live).
		StripePrice(id)

	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
