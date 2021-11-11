package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"log"
)

const (
	basePathEmailAuth        = "/auth/email"
	basePathMobileAuth       = "/auth/mobile"
	basePathPwReset          = "/auth/password-reset"
	basePathWxAuth           = "/auth/wx"
	basePathOAuth            = "/oauth/callback"
	pathEmailExists          = basePathEmailAuth + "/exists"
	pathEmailLogin           = basePathEmailAuth + "/login"
	pathEmailSignUp          = basePathEmailAuth + "/signup"
	pathEmailVerification    = basePathEmailAuth + "/verification/"
	pathMobileRequestSMS     = basePathMobileAuth + "/verification"
	pathMobileLinkEmail      = basePathMobileAuth + "/link"
	pathMobileSignUp         = basePathMobileAuth + "/signup"
	pathPwResetRequestLetter = basePathPwReset + "/letter"
	pathPwResetVerifyToken   = basePathPwReset + "/tokens/"
	pathWxLogin              = basePathWxAuth + "/login"
	pathWxRefresh            = basePathWxAuth + "/refresh"
	pathWxCallback           = basePathWxAuth + "/wx/fta-reader"
)

const (
	pathBaseAccount     = "/account"
	pathEmail           = pathBaseAccount + "/email"
	pathRequestVrfEmail = pathBaseAccount + "/email/request-verification"
	pathUserName        = pathBaseAccount + "/name"
	pathPassword        = pathBaseAccount + "/password"
	pathMobile          = pathBaseAccount + "/mobile"
	pathMobileUpdateSMS = pathBaseAccount + "/mobile/verification"
	pathAddress         = pathBaseAccount + "/address"
	pathProfile         = pathBaseAccount + "/profile"
	pathWxAccount       = pathBaseAccount + "/wx"
	pathWxSignUp        = pathBaseAccount + "/wx/signup"
	pathWxLink          = pathBaseAccount + "/wx/link"
	pathWxUnlink        = pathBaseAccount + "/wx/unlink"
)

const (
	pathBaseWxPay     = "/wxpay"
	pathBaseAliPay    = "/alipay"
	pathWxPayDesktop  = pathBaseWxPay + "/desktop"
	pathWxPayMobile   = pathBaseWxPay + "/mobile" // Mobile browser
	pathWxPayJsApi    = pathBaseWxPay + "/jsapi"  // wechat in-house browser.
	pathAliPayDesktop = pathBaseAliPay + "/desktop"
	pathAliPayMobile  = pathBaseAliPay + "/mobile"
)

type Client struct {
	key                string
	baseURL            string // Localhost for dev; otherwise online production.
	wxRedirectBaseURLs map[bool]string
	isProd             bool
}

func NewSubsAPIClient(prod bool) Client {
	log.Printf("Client for subscription api running in production: %t", prod)
	return Client{
		key:     config.MustSubsAPIKey().Pick(prod),        // Pick the correct api access token
		baseURL: config.MustSubsAPIv3BaseURL().Pick(false), // Always use localhost since this app is on the same server as API.
		wxRedirectBaseURLs: map[bool]string{
			// When this app is in production mode, we want wechat to
			// redirect to production api.
			true: "https://www.ftacademy.cn/api/v3",
			// When this app is in development mode, we want Wechat to
			// redirect to sandbox api so that changes won't affect current
			// user.
			false: "https://www.ftacademy.cn/api/sandbox",
		},
		isProd: prod,
	}
}

func (c Client) WxOAuthSession(appID string) (WxOAuthCodeSession, error) {
	midwayBaseURL := c.wxRedirectBaseURLs[c.isProd]
	cbURL := midwayBaseURL + pathWxCallback

	req := NewWxOAuthCodeRequest(appID, cbURL)

	return NewWxOAuthSession(req)
}
