package reader

import (
	"github.com/FTChinese/go-rest/rand"
	"github.com/gorilla/schema"
	"net/url"
	"time"
)

const wxOAuthURL = "https://open.weixin.qq.com/connect/qrconnect"

// The callback url must not be any sub-domains.
const callback = "https://www.ftacademy.cn/reader/wx/oauth/callback"

var encoder = schema.NewEncoder()

// wxOAuthCodeParams is used to build query parameters to get an OAuth code
// from wechat api as described here: https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
type wxOAuthCodeParams struct {
	AppID        string `schema:"appid"`
	CallbackURI  string `schema:"redirect_uri"`
	ResponseType string `schema:"response_type"`
	Scope        string `schema:"scope"`
	State        string `schema:"state"`
	Fragment     string `schema:"-"`
}

func newWxOAuthCodeParams(appID string) wxOAuthCodeParams {
	return wxOAuthCodeParams{
		AppID:        appID,
		CallbackURI:  callback,
		ResponseType: "code",
		Scope:        "snsapi_login",
		State:        rand.String(12),
		Fragment:     "wechat_redirect",
	}
}

// encodeQuery turns the struct into url-encoded query parameter.
func (p wxOAuthCodeParams) encodeQuery() (string, error) {
	var v = url.Values{}

	err := encoder.Encode(p, v)
	if err != nil {
		return "", err
	}

	return v.Encode(), nil
}

func (p wxOAuthCodeParams) url() (*url.URL, error) {
	parsed, err := url.Parse(wxOAuthURL)
	if err != nil {
		return nil, err
	}

	q, err := p.encodeQuery()
	if err != nil {
		return nil, err
	}

	parsed.RawQuery = q
	parsed.Fragment = p.Fragment

	return parsed, nil
}

// WxOAuthSession contains the data for a session of Wechat OAuth.
// You should save it on the client side so that you could perform
// verification after OAuth redirect.
type WxOAuthSession struct {
	State      string `json:"state"`     // The state contains the URL that you should verify after redirect.
	ExpiresAt  int64  `json:"expiresAt"` // Expires in 5 minutes from now.
	RedirectTo string `json:"redirectTo"`
}

func NewWxOAuthSession(appID string) (WxOAuthSession, error) {
	params := newWxOAuthCodeParams(appID)

	redirectTo, err := params.url()
	if err != nil {
		return WxOAuthSession{}, err
	}

	return WxOAuthSession{
		State:      params.State,
		ExpiresAt:  time.Now().Unix() + 5*60,
		RedirectTo: redirectTo.String(),
	}, nil
}
