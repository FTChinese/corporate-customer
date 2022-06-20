package reader

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/FTChinese/ftacademy/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PubKey struct {
	Key string `json:"key"`
}

type StripeRouter struct {
	keyHolder config.KeyHolder
	clients   api.Clients
}

func NewStripeRouter(clients api.Clients, prod bool) StripeRouter {
	return StripeRouter{
		keyHolder: config.MustStripePubKey().KeyHolder(prod),
		clients:   clients,
	}
}

func (router StripeRouter) PublishableKey(c echo.Context) error {
	live := xhttp.GetQueryLive(c)

	return c.JSON(
		http.StatusOK,
		PubKey{
			Key: router.keyHolder.Select(live),
		},
	)
}
