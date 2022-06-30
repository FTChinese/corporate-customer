package api

const (
	basePathEmailAuth        = "/auth/email"
	basePathMobileAuth       = "/auth/mobile"
	basePathPwReset          = "/auth/password-reset"
	basePathWxAuth           = "/auth/wx"
	basePathOAuth            = "/oauth/callback"
	basePathLegal            = "/legal"
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
	pathWxCallback           = basePathOAuth + "/wx/fta-reader"
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
	pathPaywall       = "/paywall"
)

const (
	pathBaseOrder = "/orders"
)

const (
	pathBaseMember  = "/membership"
	pathMemberAddOn = pathBaseMember + "/addons"
)

const (
	pathBaseApple = "/apple"
)

func pathAppleSubOf(id string) string {
	return pathBaseApple + "/" + id
}

const (
	pathBaseStripe          = "/stripe"
	pathStripePrices        = pathBaseStripe + "/prices"
	pathStripeCustomer      = pathBaseStripe + "/customers"
	pathStripeSubs          = pathBaseStripe + "/subs"
	pathStripePaymentMethod = pathBaseStripe + "/payment-methods"
	pathStripeSetupIntent   = pathBaseStripe + "/setup-intents"
)

func pathVerifyOrder(id string) string {
	return pathBaseOrder + "/" + id + "/verify-payment"
}
