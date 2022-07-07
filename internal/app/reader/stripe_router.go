package reader

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type StripeRouter struct {
	clients    api.Clients
	logger     *zap.Logger
	production bool
	keyStore   config.StripeKeyStore
}

func NewStripeRouter(
	clients api.Clients,
	production bool,
	logger *zap.Logger,
) StripeRouter {
	keyStore := config.NewStripePubKeys()
	return StripeRouter{
		clients:    clients,
		logger:     logger,
		production: production,
		keyStore:   keyStore,
	}
}

func (router StripeRouter) PublishableKey(c echo.Context) error {
	live := xhttp.GetQueryLive(c)
	// For development environment, live mode is always false.
	if !router.production {
		live = false
	}

	return c.JSON(
		http.StatusOK,
		router.keyStore.Select(live),
	)
}
