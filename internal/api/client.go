package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
)

type Client struct {
	key                string
	baseURL            string // Localhost for dev; otherwise online production.
	wxRedirectBaseURLs map[bool]string
	wxRedirectBaseURL  string
	isProd             bool // Used to determine which wx oauth redirect url to use.
}

type Clients struct {
	Sandbox Client // For test accounts
	Live    Client // For non-test accounts
}

// NewClients creates both sandbox and live clients. Then are selected at runtime.
// There are 4 combinations of live + production:
//                 Live
//            /           \
//        True            False
//      /      \       /        \
// Production  Dev   Sandbox    Dev
//
// Usually the Dev points to localhost so that in development you are not bothered with
// which account you are using.
// Among them, production/development is determined upon app launch while live/sandbox
// is determined at runtime.
func NewClients(prod bool) Clients {
	token := config.MustSubsAPIKey().Pick(prod)
	// There are 3 urls to select:
	// - Development always uses localhost
	// - Production has a live mode and test mode, determined by current user account.
	sandboxBaseURLs := config.MustSandboxAPIURL()
	liveBaseURLs := config.MustProdAPIv6BaseURL()

	return Clients{
		Sandbox: Client{
			key:     token,
			baseURL: sandboxBaseURLs.Pick(prod), // Pick localhost or production test
			// When this app is in development mode, we want Wechat to
			// redirect to sandbox api so that changes won't affect current
			// user.
			wxRedirectBaseURL: sandboxBaseURLs.Pick(true),
			isProd:            prod,
		},
		Live: Client{
			key:     token,
			baseURL: liveBaseURLs.Pick(prod), // Pick localhost or production live.
			// When this app is in production mode, we want wechat to
			// redirect to production api.
			wxRedirectBaseURL: liveBaseURLs.Pick(true),
			isProd:            prod,
		},
	}
}

func (c Clients) Select(live bool) Client {
	if live {
		return c.Live
	}

	return c.Sandbox
}

func (c Client) WxOAuthSession(appID string) (WxOAuthCodeSession, error) {
	cbURL := c.wxRedirectBaseURL + pathWxCallback

	req := NewWxOAuthCodeRequest(appID, cbURL)

	return NewWxOAuthSession(req)
}
