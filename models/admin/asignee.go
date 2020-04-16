package admin

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/guregu/null"
	"strings"
)

// Assignee represents a reader who can accept
// an invitation email, and who can be granted
// a licence.
type Assignee struct {
	FtcID    null.String `json:"ftcId" db:"ftc_id"`
	Email    null.String `json:"email" db:"email"`
	UserName null.String `json:"userName" db:"user_name"`
	IsVIP    bool        `json:"isVip" db:"is_vip"`
}

// NormalizeName tries to find a proper way to greet user
// in email.
func (a Assignee) NormalizeName() string {
	if a.UserName.Valid {
		return a.UserName.String
	}

	return strings.Split(a.Email.String, "@")[0]
}

func (a Assignee) InvitationLetter(licence Licence) (postoffice.Parcel, error) {

	data := struct {
		Name string
		Letter
	}{
		Name: a.NormalizeName(),
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
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网B2B]会员邀请",
		Body:        body.String(),
	}, nil
}

type AssigneeSchema struct {
	Assignee
	TeamID string `db:"team_id"`
}
