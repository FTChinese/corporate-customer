package admin

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"strings"
)

// BaseAccount contains the shared fields of an admin.
type BaseAccount struct {
	ID          string      `json:"id" db:"admin_id"`
	TeamID      null.String `json:"teamId" db:"team_id"`
	Email       string      `json:"email" db:"email"`
	DisplayName null.String `json:"displayName" db:"display_name"`
	Active      bool        `json:"active" db:"active"`
	Verified    bool        `json:"verified" db:"verified"`
}

func (a BaseAccount) NormalizeName() string {
	if a.DisplayName.Valid {
		return a.DisplayName.String
	}

	return strings.Split(a.Email, "@")[0]
}

func (a BaseAccount) UpdateName(name string) BaseAccount {
	a.DisplayName = null.StringFrom(name)
	return a
}

// Account is used to created an admin.
type Account struct {
	BaseAccount
	Password   string      `json:"-" db:"password"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

// NewAccount creates a new instance from signup parameters.
func NewAccount(p input.SignupParams) Account {
	return Account{
		BaseAccount: BaseAccount{
			ID:          uuid.New().String(),
			TeamID:      null.String{},
			Email:       p.Email,
			DisplayName: null.String{},
			Active:      true,
			Verified:    false,
		},
		Password:   p.Password,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}
