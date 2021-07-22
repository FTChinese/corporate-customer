package config

const (
	SiteBaseURL           = "https://next.ftacademy.cn"
	B2BBaseURL            = SiteBaseURL + "/corporate"
	B2BReaderVerification = B2BBaseURL + "/user/verification"
)

func B2BVerifyInvitationURL(token string) string {
	return B2BBaseURL + "/verify-invitation/" + token
}

func B2BVerifyAdminURL(token string) string {
	return B2BBaseURL + "/verify/" + token
}
