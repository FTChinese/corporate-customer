package api

import "github.com/FTChinese/ftacademy/pkg/fetch"

func (c Client) LoadAccountByFtcID(id string) (fetch.Response, error) {
	url := c.baseURL + pathBaseAccount

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		SetFtcID(id).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

func (c Client) LoadAccountByUnionID(id string) (fetch.Response, error) {
	url := c.baseURL + pathWxAccount

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		SetUnionID(id).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}
