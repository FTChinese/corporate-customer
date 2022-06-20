package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
)

type Client struct {
	key                string
	baseURL            string // Localhost for dev; otherwise online production.
	wxRedirectBaseURLs map[bool]string
	isProd             bool
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
	sandboxBaseURLs := config.MustAPISandboxURL()
	liveBaseURLs := config.MustSubsAPIv6BaseURL()

	// To relay wechat OAuth, localhost should never be used.
	wxRelayBaseURLs := map[bool]string{
		// When this app is in production mode, we want wechat to
		// redirect to production api.
		true: liveBaseURLs.Pick(true),
		// When this app is in development mode, we want Wechat to
		// redirect to sandbox api so that changes won't affect current
		// user.
		false: sandboxBaseURLs.Pick(true),
	}

	return Clients{
		Sandbox: Client{
			key:                token,
			baseURL:            sandboxBaseURLs.Pick(prod), // Pick localhost or production test
			wxRedirectBaseURLs: wxRelayBaseURLs,
			isProd:             prod,
		},
		Live: Client{
			key:                token,
			baseURL:            liveBaseURLs.Pick(prod), // Pick localhost or production live.
			wxRedirectBaseURLs: wxRelayBaseURLs,
			isProd:             prod,
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
	midwayBaseURL := c.wxRedirectBaseURLs[c.isProd]
	cbURL := midwayBaseURL + pathWxCallback

	req := NewWxOAuthCodeRequest(appID, cbURL)

	return NewWxOAuthSession(req)
}
