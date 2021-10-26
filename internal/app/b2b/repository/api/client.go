package api

import (
	"github.com/FTChinese/ftacademy/pkg/config"
	"log"
)

const (
	pathBaseEmailAuth        = "/auth/email"
	pathBaseMobileAuth       = "/auth/mobile"
	pathBasePwReset          = "/password-reset"
	pathBaseAccount          = "/account"
	pathEmailExists          = pathBaseEmailAuth + "/exists"
	pathEmailLogin           = pathBaseEmailAuth + "/login"
	pathEmailSignUp          = pathBaseEmailAuth + "/signup"
	pathEmailVerification    = pathBaseEmailAuth + "/verification/"
	pathMobileRequestSMS     = pathBaseMobileAuth + "/verification"
	pathMobileLinkEmail      = pathBaseMobileAuth + "/link"
	pathMobileSignUp         = pathBaseMobileAuth + "/signup"
	pathPwResetRequestLetter = pathBasePwReset + "/letter"
	pathPwResetVerifyToken   = pathBasePwReset + "/tokens/"
	pathEmail                = pathBaseAccount + "/email"
	pathRequestVrfEmail      = pathBaseAccount + "/email/request-verification"
	pathUserName             = pathBaseAccount + "/password"
	pathMobile               = pathBaseAccount + "/mobile"
	pathSMSNewMobile         = pathBaseAccount + "/mobile/verification"
	pathAddress              = pathBaseAccount + "/address"
	pathProfile              = pathBaseAccount + "/profile"
	pathWxAccount            = pathBaseAccount + "/wx"
	pathWxSignUp             = pathBaseAccount + "/wx/signup"
	pathWxLink               = pathBaseAccount + "/wx/link"
	pathWxUnlink             = pathBaseAccount + "/wx/unlink"
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
