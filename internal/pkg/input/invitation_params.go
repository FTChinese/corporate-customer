package input

import (
	"github.com/FTChinese/ftacademy/pkg/validator"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"strings"
)

// InvitationParams contains the essential data client
// submitted to create a new invitation.
type InvitationParams struct {
	Email       string      `json:"email"` // To whom the invitation should be sent.
	Description null.String `json:"description"`
	LicenceID   string      `json:"licenceId"` // Which licence is being granted.
	TeamID      string      `json:"teamId"`
}

func (i *InvitationParams) Validate() *render.ValidationError {
	i.Email = strings.TrimSpace(i.Email)
	desc := strings.TrimSpace(i.Description.String)
	i.Description = null.NewString(desc, desc != "")
	i.LicenceID = strings.TrimSpace(i.LicenceID)
	i.TeamID = strings.TrimSpace(i.TeamID)

	ve := validator.New("email").Required().Email().Validate(i.Email)
	if ve != nil {
		return ve
	}

	ve = validator.New("description").MaxLen(128).Validate(i.Description.String)
	if ve != nil {
		return ve
	}

	ve = validator.New("licenceId").Required().Validate(i.LicenceID)
	if ve != nil {
		return ve
	}

	return validator.New("teamId").Required().Validate(i.TeamID)
}
