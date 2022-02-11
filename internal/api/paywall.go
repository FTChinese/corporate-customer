package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) Paywall() (*http.Response, error) {
	url := c.baseURL + pathPaywall

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ListStripePrices() (*http.Response, error) {
	url := c.baseURL + pathStripePrices

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripePrice(id string) (*http.Response, error) {
	url := c.baseURL + pathStripePrices + "/" + id

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
