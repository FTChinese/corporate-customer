package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"net/http"
)

func (c Client) ResetPassword(body io.Reader) (*http.Response, error) {
	url := c.baseURL + basePathPwReset

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) RequestPasswordResetLetter(body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathPwResetRequestLetter

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) VerifyResetToken(token string) (*http.Response, error) {
	url := c.baseURL + pathPwResetVerifyToken + token

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
