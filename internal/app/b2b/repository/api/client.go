package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"log"
)

const (
	pathBaseEmailAuth        = "/auth/email"
	pathBaseMobileAuth       = "/auth/mobile"
	pathBasePwReset          = "/auth/password-reset"
	pathEmailExists          = pathBaseEmailAuth + "/exists"
	pathEmailLogin           = pathBaseEmailAuth + "/login"
	pathEmailSignUp          = pathBaseEmailAuth + "/signup"
	pathEmailVerification    = pathBaseEmailAuth + "/verification/"
	pathMobileRequestSMS     = pathBaseMobileAuth + "/verification"
	pathMobileLinkEmail      = pathBaseMobileAuth + "/link"
	pathMobileSignUp         = pathBaseMobileAuth + "/signup"
	pathPwResetRequestLetter = pathBasePwReset + "/letter"
	pathPwResetVerifyToken   = pathBasePwReset + "/tokens/"
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
	key     string
	baseURL string
}

func NewSubsAPIClient(prod bool) Client {
	log.Printf("Client for subscription api running in production: %t", prod)
	return Client{
		key:     config.MustSubsAPIKey().Pick(prod),        // Pick the correct api access token
		baseURL: config.MustSubsAPIv3BaseURL().Pick(false), // Always use localhost.
	}
}
