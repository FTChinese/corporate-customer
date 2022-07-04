package reader

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type StripeRouter struct {
	keyHolder config.StripeKeyHolder
	clients   api.Clients
	logger    *zap.Logger
}

func NewStripeRouter(
	clients api.Clients,
	keyHolder config.StripeKeyHolder,
	logger *zap.Logger,
) StripeRouter {
	return StripeRouter{
		keyHolder: keyHolder,
		clients:   clients,
		logger:    logger,
	}
}

func (router StripeRouter) PublishableKey(c echo.Context) error {

	return c.JSON(
		http.StatusOK,
		router.keyHolder,
	)
}
