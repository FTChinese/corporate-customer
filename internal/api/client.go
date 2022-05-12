package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"log"
)

type Client struct {
	key                string
	baseURL            string // Localhost for dev; otherwise online production.
	wxRedirectBaseURLs map[bool]string
	isProd             bool
}

func NewSubsAPIClient(prod bool) Client {
	log.Printf("Client for subscription api running in production: %t", prod)

	baseURLs := config.MustSubsAPIv6BaseURL()
	return Client{
		key:     config.MustSubsAPIKey().Pick(prod), // Pick the correct api access token
		baseURL: baseURLs.Pick(false),               // Always use localhost since this app is on the same server as API.
		wxRedirectBaseURLs: map[bool]string{
			// When this app is in production mode, we want wechat to
			// redirect to production api.
			true: baseURLs.Pick(true),
			// When this app is in development mode, we want Wechat to
			// redirect to sandbox api so that changes won't affect current
			// user.
			false: config.MustAPISandboxURL().Pick(true),
		},
		isProd: prod,
	}
}

type Clients struct {
	Test Client
	Live Client
}

func NewClients(prod bool) Clients {
	token := config.MustSubsAPIKey().Pick(prod)
	// There are 3 urls to select:
	// - Development always uses localhost
	// - Production has a live mode and test mode, determined by current user account.
	testURLSelector := config.MustAPISandboxURL()
	liveURLSelector := config.MustSubsAPIv6BaseURL()

	// To relay wechat OAuth, localhost should never be used.
	wxRelayBaseURLs := map[bool]string{
		// When this app is in production mode, we want wechat to
		// redirect to production api.
		true: liveURLSelector.Pick(true),
		// When this app is in development mode, we want Wechat to
		// redirect to sandbox api so that changes won't affect current
		// user.
		false: testURLSelector.Pick(true),
	}

	return Clients{
		Test: Client{
			key:                token,
			baseURL:            config.MustAPISandboxURL().Pick(prod), // Pick localhost or production test
			wxRedirectBaseURLs: wxRelayBaseURLs,
			isProd:             prod,
		},
		Live: Client{
			key:                token,
			baseURL:            config.MustSubsAPIv6BaseURL().Pick(prod), // Pick localhost or production live.
			wxRedirectBaseURLs: wxRelayBaseURLs,
			isProd:             prod,
		},
	}
}

func (c Clients) Select(live bool) Client {
	if live {
		return c.Live
	}

	return c.Test
}

func (c Client) WxOAuthSession(appID string) (WxOAuthCodeSession, error) {
	midwayBaseURL := c.wxRedirectBaseURLs[c.isProd]
	cbURL := midwayBaseURL + pathWxCallback

	req := NewWxOAuthCodeRequest(appID, cbURL)

	return NewWxOAuthSession(req)
}
