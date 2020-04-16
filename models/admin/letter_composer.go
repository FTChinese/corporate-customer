package admin

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/postoffice"
	"strings"
)

const baseUrl = "https://www.ftacademy.cn/b2b"

// Letter is the data passed to template to generate the content
// of an email.
type Letter struct {
	URL      string // The link for verification, or password reset.
	IsSignUp bool   // Determines greeting.
}

type InvitationLetter struct {
	AdminEmail string
	TeamName   string
	Tier       enum.Tier
	URL        string
}

func ComposeInvitationLetter(a Assignee, l Licence, at AccountTeam) (postoffice.Parcel, error) {
	data := struct {
		AssigneeName string
		TeamName     string
		Tier         enum.Tier
		URL          string
		AdminEmail   string
	}{
		AssigneeName: a.NormalizeName(),
		TeamName:     at.TeamName.String,
		Tier:         l.Plan.Tier,
		URL:          baseUrl + "/accept-invitation/" + l.Invitation.Token,
		AdminEmail:   at.Email,
	}

	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "invitation", data)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email.String,
		ToName:      data.AssigneeName,
		Subject:     "[FT中文网B2B]会员邀请",
		Body:        body.String(),
	}, nil
}
