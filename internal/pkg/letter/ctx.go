package letter

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

const baseUrl = "https://next.ftacademy.cn/b2b"

type CtxVerification struct {
	Email    string
	UserName string
	Link     string
	IsSignUp bool
}

func (ctx CtxVerification) Render() (string, error) {
	return Render(keyVrf, ctx)
}

type CtxVerified struct {
	UserName string
}

func (ctx CtxVerified) Render() (string, error) {
	return Render(keyVerified, ctx)
}

type CtxPwReset struct {
	UserName string
	Link     string
	Duration string
}

func (ctx CtxPwReset) Render() (string, error) {
	return Render(keyPwReset, ctx)
}

type CtxInvitation struct {
	AssigneeName string
	AdminEmail   string
	TeamName     string
	Tier         enum.Tier
	URL          string
}

func (ctx CtxInvitation) Render() (string, error) {
	return Render(keyLicenceInvitation, ctx)
}

type CtxLicenceGranted struct {
	Name           string
	AssigneeEmail  string
	Tier           enum.Tier
	ExpirationDate chrono.Date
}

func (ctx CtxLicenceGranted) Render() (string, error) {
	return Render(keyLicenceGranted, ctx)
}
