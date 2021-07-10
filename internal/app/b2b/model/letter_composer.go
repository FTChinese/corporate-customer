package model

import (
	"github.com/FTChinese/ftacademy/pkg/postman"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"strings"
)

const baseUrl = "https://next.ftacademy.cn/b2b"

type InvitationLetter struct {
	AdminEmail string
	TeamName   string
	Tier       enum.Tier
	URL        string
}

func ComposeInvitationLetter(il InvitedLicence, pp Passport) (postman.Parcel, error) {
	data := struct {
		AssigneeName string
		TeamName     string
		Tier         enum.Tier
		URL          string
		AdminEmail   string
	}{
		AssigneeName: il.Assignee.NormalizeName(),
		TeamName:     pp.TeamName.String,
		Tier:         il.Plan.Tier,
		URL:          baseUrl + "/accept-invitation/" + il.Invitation.Token,
		AdminEmail:   pp.Email,
	}

	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "invitation", data)

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   il.Assignee.Email.String,
		ToName:      data.AssigneeName,
		Subject:     "[FT中文网B2B]会员邀请",
		Body:        body.String(),
	}, nil
}

func ComposeVerificationLetter(a Account, verifier AccountInput) (postman.Parcel, error) {

	data := struct {
		Name     string
		URL      string
		IsSignUp bool
	}{
		Name:     a.NormalizeName(),
		URL:      baseUrl + "/verify/" + verifier.Token,
		IsSignUp: verifier.IsSignUp,
	}
	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "verification", data)

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网B2B]验证账号",
		Body:        body.String(),
	}, nil
}

func ComposePwResetLetter(a Account, bearer AccountInput) (postman.Parcel, error) {

	data := struct {
		Name string
		URL  string
	}{
		Name: a.NormalizeName(),
		URL:  baseUrl + "/password-reset/token/" + bearer.Token,
	}
	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "passwordReset", data)

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      data.Name,
		Subject:     "[FT中文网B2B]重置密码",
		Body:        body.String(),
	}, nil
}

// ComposeLicenceGranted write a letter to admin after
// a reader accepted an invitation and the corresponding
// licence is granted.
// We need to know the admin's account, reader's email
// the the licence's plan.
func ComposeLicenceGranted(il InvitedLicence, pp Passport) (postman.Parcel, error) {

	var data = struct {
		Name           string
		AssigneeEmail  string
		Tier           enum.Tier
		ExpirationDate chrono.Date
	}{
		Name:           pp.NormalizeName(),
		AssigneeEmail:  il.Assignee.Email.String,
		Tier:           il.Plan.Tier,
		ExpirationDate: il.Licence.ExpireDate,
	}

	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "licenceGranted", data)

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   pp.Email,
		ToName:      data.Name,
		Subject:     "[FT中文网B2B]团队成员获得会员许可",
		Body:        body.String(),
	}, nil
}
