package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"net/http"
)

func (c Client) Membership(ids reader.PassportClaims) (*http.Response, error) {
	u := c.baseURL + pathBaseMember

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

func (c Client) ClaimAddOn(ids reader.PassportClaims) (*http.Response, error) {
	url := c.baseURL + pathMemberAddOn

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
