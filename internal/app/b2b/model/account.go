package model

import (
	"github.com/FTChinese/go-rest/chrono"
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

type Profile struct {
	Account
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}
