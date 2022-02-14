package reader

import (
	"github.com/FTChinese/ftacademy/internal/api"
	"github.com/FTChinese/ftacademy/pkg/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PubKey struct {
	Key string `json:"key"`
}

type StripeRouter struct {
	pubKey    string
	apiClient api.Client
}

func NewStripeRouter(client api.Client, prod bool) StripeRouter {
	return StripeRouter{
		pubKey:    config.MustStripePubKey().Pick(prod),
		apiClient: client,
	}
}

func (router StripeRouter) PublishableKey(c echo.Context) error {
	return c.JSON(http.StatusOK, PubKey{Key: router.pubKey})
}
