package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"log"
	"net/http"
)

func (c Client) StripeNewCustomer(ids reader.PassportClaims) (*http.Response, error) {
	url := c.baseURL + pathStripeCustomer

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetCustomer(ids reader.PassportClaims, cusID string) (*http.Response, error) {
	url := c.baseURL + pathCustomerOf(cusID)

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeCustomerDefaultPaymentMethod(ids reader.PassportClaims, cusID string) (*http.Response, error) {
	url := c.baseURL + pathCustomerOf(cusID) + "/default-payment-method"

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeNewSubs(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathStripeSubs

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSubs(ids reader.PassportClaims) (*http.Response, error) {
	url := c.baseURL + pathStripeSubs

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id)

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeUpdateSubsOf(ids reader.PassportClaims, id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id)

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeRefreshSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id) + "/refresh"

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeCancelSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id) + "/cancel"

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeReactivateSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id) + "/reactivate"

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeSubsDefaultPaymentMethod(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathSubsOf(id) + "/default-payment-method"

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeListPaymentMethods(ids reader.PassportClaims, rawQuery string) (*http.Response, error) {
	url := c.baseURL + pathStripePaymentMethod + "?" + rawQuery

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripePaymentMethodOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	url := c.baseURL + pathPaymentMethodOf(id)

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
