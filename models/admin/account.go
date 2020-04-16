package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/guregu/null"
	"strings"
)

// Account is an organization's administrator account.
// An account might manage multiple teams/organizations.
// Currently we allow only one team per account.
type Account struct {
	ID          string      `json:"id" db:"admin_id"`
	Email       string      `json:"email" db:"email"`
	DisplayName null.String `json:"displayName" db:"display_name"`
	Active      bool        `json:"active" db:"active"`
	Verified    bool        `json:"verified" db:"verified"`
}

func (a Account) NormalizeName() string {
	if a.DisplayName.Valid {
		return a.DisplayName.String
	}

	return strings.Split(a.Email, "@")[0]
}

func (a Account) VerificationLetter(letter Letter) (postoffice.Parcel, error) {

	data := struct {
		Name string
		Letter
	}{
		Name:   a.NormalizeName(),
		Letter: letter,
	}
	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "verification", data)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网B2B]验证账号",
		Body:        body.String(),
	}, nil
}

func (a Account) PasswordResetLetter(letter Letter) (postoffice.Parcel, error) {

	data := struct {
		Name string
		Letter
	}{
		Name:   a.NormalizeName(),
		Letter: letter,
	}
	var body strings.Builder
	err := tmpl.ExecuteTemplate(&body, "passwordReset", data)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      data.Name,
		Subject:     "[FT中文网B2B]重置密码",
		Body:        body.String(),
	}, nil
}

type Profile struct {
	Account
	CreatedUTC chrono.Time `db:"created_utc"`
	UpdatedUTC chrono.Time `db:"updated_utc"`
}

type AccountTeam struct {
	Account
	TeamID   null.String `json:"teamId" db:"team_id"`
	TeamName null.String `json:"-" db:"team_name"`
}
