package letter

// CtxVerification holds data to render a letter upon signup.
type CtxVerification struct {
	Email    string
	UserName string
	Link     string
}

func (ctx CtxVerification) Render() (string, error) {
	return Render(keyVrf, ctx)
}

// CtxVerified is used to compose an email after email verified.
type CtxVerified struct {
	UserName string
}

func (ctx CtxVerified) Render() (string, error) {
	return Render(keyVerified, ctx)
}

// CtxPwReset is used to compose an email upon requesting
// password reset.
type CtxPwReset struct {
	UserName string
	Link     string
	Duration string
}

func (ctx CtxPwReset) Render() (string, error) {
	return Render(keyPwReset, ctx)
}

// CtxInvitation is used to compose an invitation email
// so that B2B org's member could use a licence.
type CtxInvitation struct {
	ToName     string
	AdminEmail string
	TeamName   string
	Tier       string
	URL        string
}

func (ctx CtxInvitation) Render() (string, error) {
	return Render(keyLicenceInvitation, ctx)
}

// CtxLicenceGranted is used to notify admin that a licence
// is granted to a member.
type CtxLicenceGranted struct {
	Name           string
	AssigneeEmail  string
	Tier           string
	ExpirationDate string
}

func (ctx CtxLicenceGranted) Render() (string, error) {
	return Render(keyLicenceGranted, ctx)
}
