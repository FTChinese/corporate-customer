package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) Paywall() (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathPaywall).
		String()

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
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
		String()

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
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
		AddPath(id).
		String()

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
