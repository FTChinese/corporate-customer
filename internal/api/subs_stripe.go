package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"log"
	"net/http"
)

func (c Client) StripeNewCustomer(ids reader.PassportClaims) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeCustomer).String()

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

	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeCustomer).AddPath(cusID).String()

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

func (c Client) StripeCusDefaultPaymentMethod(ids reader.PassportClaims, cusID string) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeCustomer).
		AddPath(cusID).
		AddPath("default-payment-method").
		String()

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

func (c Client) StripeSetCusDefaultPaymentMethod(ids reader.PassportClaims, cusID string, body io.Reader) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeCustomer).
		AddPath(cusID).
		AddPath("default-payment-method").
		String()

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

func (c Client) StripeListCusPaymentMethods(ids reader.PassportClaims, cusID string, rawQuery string) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeCustomer).
		AddPath(cusID).
		AddPath("payment-methods").
		SetRawQuery(rawQuery).
		String()

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

func (c Client) StripeCreateSetupIntent(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeSetupIntent).String()

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

func (c Client) StripeGetSetupIntent(ids reader.PassportClaims, id string, rawQuery string) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSetupIntent).
		AddPath(id).
		SetRawQuery(rawQuery).
		String()

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

func (c Client) StripeGetSetupPaymentMethod(ids reader.PassportClaims, id string, rawQuery string) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSetupIntent).
		AddPath(id).
		AddPath("payment-method").
		SetRawQuery(rawQuery).
		String()

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
	u := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeSubs).String()

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
	u := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeSubs).AddPath(id).String()

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
	u := fetch.NewURLBuilder(c.baseURL).AddPath(pathStripeSubs).AddPath(id).String()

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
	u := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(id).
		AddPath("refresh").
		String()

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
	u := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(id).
		AddPath("cancel").
		String()

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
	u := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(id).
		AddPath("reactivate").
		String()

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
	u := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(id).
		AddPath("default-payment-method").
		String()

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

func (c Client) StripeSetSubsDefaultPaymentMethod(ids reader.PassportClaims, subsID string, body io.Reader) (*http.Response, error) {

	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(subsID).
		AddPath("default-payment-method").
		String()

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

// StripeSubsLatestInvoice loads a subscription's latest invoice.
func (c Client) StripeSubsLatestInvoice(claims reader.PassportClaims, subsID string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(subsID).
		AddPath("latest-invoice").
		String()

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(claims).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CouponOfLatestSubsInvoice(claims reader.PassportClaims, subsID string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeSubs).
		AddPath(subsID).
		AddPath("latest-invoice/any-coupon").
		String()

	log.Printf("Fetching data from %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		WithHeader(ReaderIDsHeader(claims).Build()).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) StripePaymentMethodOf(ids reader.PassportClaims, id string, rawQuery string) (*http.Response, error) {
	u := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePaymentMethod).
		AddPath(id).
		SetRawQuery(rawQuery).
		String()

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
