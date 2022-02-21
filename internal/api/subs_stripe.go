package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (c Client) StripeNewCustomer(ids reader.PassportClaims) (*http.Response, error) {
	u := c.baseURL + pathStripeCustomer

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetCustomer(ids reader.PassportClaims, cusID string) (*http.Response, error) {
	u := c.baseURL + pathCustomerOf(cusID)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeCusDefaultPaymentMethod(ids reader.PassportClaims, cusID string) (*http.Response, error) {
	u := c.baseURL + pathCusDefaultPaymentMethod(cusID)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeSetCusDefaultPaymentMethod(ids reader.PassportClaims, cusID string, body io.Reader) (*http.Response, error) {
	u := c.baseURL + pathCusDefaultPaymentMethod(cusID)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeListCusPaymentMethods(ids reader.PassportClaims, cusID string, q url.Values) (*http.Response, error) {
	u := c.baseURL + pathCusPaymentMethods(cusID)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		WithQuery(q).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeCreateSetupIntent(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	u := c.baseURL + pathStripeSetupIntent

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSetupIntent(ids reader.PassportClaims, id string, q url.Values) (*http.Response, error) {
	u := c.baseURL + pathSetupIntentOf(id)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		WithQuery(q).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSetupPaymentMethod(ids reader.PassportClaims, id string, q url.Values) (*http.Response, error) {
	u := c.baseURL + pathSetupIntentOf(id) + "/payment-method"

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		WithQuery(q).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeNewSubs(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	u := c.baseURL + pathStripeSubs

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSubs(ids reader.PassportClaims, id string) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeGetSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeUpdateSubsOf(ids reader.PassportClaims, id string, body io.Reader) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
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
	u := c.baseURL + pathSubsOf(id) + "/refresh"

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeCancelSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id) + "/cancel"

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeReactivateSubsOf(ids reader.PassportClaims, id string) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id) + "/reactivate"

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Post(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripeSubsDefaultPaymentMethod(ids reader.PassportClaims, id string) (*http.Response, error) {
	u := c.baseURL + pathSubsOf(id) + "/default-payment-method"

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripePaymentMethodOf(ids reader.PassportClaims, id string, query url.Values) (*http.Response, error) {
	u := c.baseURL + pathPaymentMethodOf(id)

	log.Printf("Fetching data from %s", u)

	resp, errs := fetch.
		New().
		Get(u).
		WithHeader(ReaderIDsHeader(ids).Build()).
		SetBearerAuth(c.key).
		WithQuery(query).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
