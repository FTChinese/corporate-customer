package api

import (
	"github.com/FTChinese/ftacademy/pkg/fetch"
	"io"
	"net/http"
)

func (c Client) LoadAccountByFtcID(id string) (fetch.Response, error) {
	url := c.baseURL + pathBaseAccount

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
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
		SetHeader(keyUnionID, id).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

// UpdateEmail changes email.
// * email: string;
// * sourceUrl: string. Used to compose verification link.
func (c Client) UpdateEmail(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathEmail

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RequestEmailVerification sends a verification letter to user's email.
func (c Client) RequestEmailVerification(id string) (*http.Response, error) {
	url := c.baseURL + pathRequestVrfEmail

	resp, errs := fetch.
		New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// UpdateName changes username.
// * userName: string;
func (c Client) UpdateName(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathUserName

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// UpdatePassword changes password.
// * currentPassword: string;
// * newPassword: string;
func (c Client) UpdatePassword(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathPassword

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RequestMobileUpdateSMS requests an SMS before switching mobile.
// * mobile: string
func (c Client) RequestMobileUpdateSMS(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathMobileUpdateSMS

	resp, errs := fetch.
		New().
		Put(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// UpdateMobile switches mobile to a new one.
// * mobile: string;
// * code: string;
func (c Client) UpdateMobile(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathMobile

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadAddress(id string) (*http.Response, error) {
	url := c.baseURL + pathAddress

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// UpdateAddress changes the following fields:
// * country?: string;
// * province?: string;
// * city?: string;
// * district?: string;
// * street?: string;
// * postcode?: string
func (c Client) UpdateAddress(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathAddress

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadProfile(id string) (*http.Response, error) {
	url := c.baseURL + pathProfile

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// UpdateProfile updates the following fields:
// * familyName?: string;
// * givenName?: string;
// * birthday?: string;
// * gender?: M | F;
func (c Client) UpdateProfile(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathProfile

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		SetHeader(keyFtcID, id).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// WxSignUp enables a wechat login user to create a
// new email account.
// * email: string;
// * password: string;
// * sourceUrl: string;
//
// The footprint.Client headers are required.
func (c Client) WxSignUp(unionID string, body io.Reader) (fetch.Response, error) {
	url := c.baseURL + pathWxSignUp

	resp, errs := fetch.
		New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(keyUnionID, unionID).
		StreamJSON(body).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

// WxLink lets a wechat-only user to link an existing
// email account.
// Header
// * `X-Union-Id: <wechat union id>`
//
// Input: input.LinkWxParams
// * ftcId: string;
// Returns no data.
func (c Client) WxLink(unionID string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathWxLink

	resp, errs := fetch.
		New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(keyUnionID, unionID).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// WxUnlink lets a wechat-only user to link an existing
// email account.
// * ftcId: string
// * anchor?: ftc | wechat;
// Returns no data.
func (c Client) WxUnlink(unionID string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathWxUnlink

	resp, errs := fetch.
		New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(keyUnionID, unionID).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
