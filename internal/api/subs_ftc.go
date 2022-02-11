package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"net/http"
)

// WxPayDesktop handles payment in desktop browsers.
// * priceId: string;
// * discountId: string;
func (c Client) WxPayDesktop(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathWxPayDesktop

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

// WxPayMobile handles payment in desktop browsers.
// * priceId: string;
// * discountId: string;
func (c Client) WxPayMobile(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathWxPayMobile

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

// WxPayJsApi handles payment in desktop browsers.
// * priceId: string;
// * discountId: string;
// * openId: string; trade_type=JSAPI时（即JSAPI支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识。
func (c Client) WxPayJsApi(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathWxPayJsApi

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

// AliPayDesktop handles payment in desktop browsers.
// * priceId: string;
// * discountId: string;
func (c Client) AliPayDesktop(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathAliPayDesktop

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

// AliPayMobile handles payment in desktop browsers.
// * priceId: string;
// * discountId: string;
func (c Client) AliPayMobile(ids reader.PassportClaims, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathAliPayMobile

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
