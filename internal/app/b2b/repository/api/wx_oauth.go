package api

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/gorilla/schema"
	"net/url"
)

var encoder = schema.NewEncoder()

// WxOAuthCodeRequest is used to build the url to get an OAuth code
// from wechat api as described here: https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
type WxOAuthCodeRequest struct {
	CodeRequestURL url.URL `schema:"-"` // The base url of wechat where we ask for user consent.
	AppID          string  `schema:"appid"`
	CallbackURI    string  `schema:"redirect_uri"`
	ResponseType   string  `schema:"response_type"`
	Scope          string  `schema:"scope"`
	State          string  `schema:"state"`
	Fragment       string  `schema:"-"`
}

// NewWxOAuthCodeRequest creates the parameters required to
// request an OAuth code.
// isProd determines where Wechat should direct us after user
// consent.
func NewWxOAuthCodeRequest(appID string, cb string) WxOAuthCodeRequest {

	return WxOAuthCodeRequest{
		CodeRequestURL: url.URL{
			Scheme:   "https",
			Host:     "open.weixin.qq.com",
			Path:     "/connect/qrconnect",
			RawQuery: "", // Query parameters to be populated.
			Fragment: "", // To be populated
		},
		AppID:        appID,
		CallbackURI:  cb,
		ResponseType: "code",
		Scope:        "snsapi_login",
		State:        rand.String(12),
		Fragment:     "wechat_redirect",
	}
}

// EncodeQuery turns the struct into url-encoded query parameter.
func (p WxOAuthCodeRequest) EncodeQuery() (string, error) {
	var v = url.Values{}

	err := encoder.Encode(p, v)
	if err != nil {
		return "", err
	}

	return v.Encode(), nil
}

// Build the complete url to request oauth code.
func (p WxOAuthCodeRequest) Build() (string, error) {
	q, err := p.EncodeQuery()
	if err != nil {
		return "", err
	}

	p.CodeRequestURL.RawQuery = q
	p.CodeRequestURL.Fragment = p.Fragment

	return p.CodeRequestURL.String(), nil
}

// WxOAuthCodeSession contains the data for a session of Wechat OAuth.
// You should save it on the client side so that you could perform
// verification after OAuth redirect.
type WxOAuthCodeSession struct {
	State      string `json:"state"` // The state contains the URL that you should verify after redirect.
	RedirectTo string `json:"redirectTo"`
}

func NewWxOAuthSession(params WxOAuthCodeRequest) (WxOAuthCodeSession, error) {

	redirectTo, err := params.Build()
	if err != nil {
		return WxOAuthCodeSession{}, err
	}

	return WxOAuthCodeSession{
		State:      params.State,
		RedirectTo: redirectTo,
	}, nil
}

// WxLoginSession is the response of wechat login endpoint.
type WxLoginSession struct {
	ID        string      `json:"sessionId"` // use this to refresh account.
	UnionID   string      `json:"unionId"`   // Use this to fetch account data
	CreatedAt chrono.Time `json:"createdAt"`
}
