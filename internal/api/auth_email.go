package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"log"
	"net/http"
)

func (c Client) EmailExists(rawQuery string) (*http.Response, error) {
	url := c.baseURL + pathEmailExists + "?" + rawQuery

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) EmailLogin(body io.Reader, client http.Header) (fetch.Response, error) {
	url := c.baseURL + pathEmailLogin

	resp, errs := fetch.
		New().
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

// EmailSignUp forwards email signup request to API and
// returns the response status code and body as bytes for
// further modification.
func (c Client) EmailSignUp(body io.Reader, client http.Header) (fetch.Response, error) {
	url := c.baseURL + pathEmailSignUp

	log.Printf("Forward reader signup to api %s", url)

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

func (c Client) VerifyEmail(token string) (*http.Response, error) {
	url := c.baseURL + pathEmailVerification + token

	log.Printf("Forward reader verification to api %s", url)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
