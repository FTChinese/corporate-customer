package input

import (
	"github.com/FTChinese/ftacademy/pkg/validator"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"strings"
)

type EmailUpdateParams struct {
	Email     string      `json:"email"`
	SourceURL null.String `json:"sourceUrl"`
}

func (i *EmailUpdateParams) Validate() *render.ValidationError {
	i.Email = strings.TrimSpace(i.Email)
	url := strings.TrimSpace(i.SourceURL.String)

	i.SourceURL = null.NewString(url, url != "")

	return validator.EnsureEmail(i.Email)
}

type ReqEmailVrfParams struct {
	SourceURL null.String `json:"sourceUrl"`
}

type NameUpdateParams struct {
	DisplayName string `json:"displayName"`
}

func (p *NameUpdateParams) Validate() *render.ValidationError {
	p.DisplayName = strings.TrimSpace(p.DisplayName)

	return validator.New("userName").
		Required().
		MaxLen(64).
		Validate(p.DisplayName)
}

type PasswordUpdateParams struct {
	ID  string `json:"-" db:"id"`
	Old string `json:"oldPassword"`
	New string `json:"password" db:"password"` // required. max 128 chars
}

// Validate is only used when a logged-in user changing password.
func (a *PasswordUpdateParams) Validate() *render.ValidationError {
	a.New = strings.TrimSpace(a.New)
	a.Old = strings.TrimSpace(a.Old)

	ve := validator.EnsurePassword(a.New)
	if ve != nil {
		return ve
	}

	return validator.
		New("oldPassword").
		Required().
		Validate(a.Old)
}
