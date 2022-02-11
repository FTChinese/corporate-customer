package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"net/http"
)

func (c Client) RequestLoginSMS(body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathMobileRequestSMS

	resp, errs := fetch.New().
		Put(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) VerifyLoginSMS(body io.Reader) (fetch.Response, error) {
	url := c.baseURL + pathMobileRequestSMS

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

func (c Client) MobileLinkExistingEmail(body io.Reader, client http.Header) (fetch.Response, error) {
	url := c.baseURL + pathMobileLinkEmail

	resp, errs := fetch.New().
		Post(url).
		WithHeader(client).
		SetBearerAuth(c.key).
		StreamJSON(body).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

func (c Client) MobileSignUp(body io.Reader, client http.Header) (fetch.Response, error) {
	url := c.baseURL + pathMobileSignUp

	resp, errs := fetch.New().
		Post(url).
		WithHeader(client).
		SetBearerAuth(c.key).
		StreamJSON(body).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}
