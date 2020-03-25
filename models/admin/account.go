package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// Account is an organization's administrator account.
// An account might manage multiple teams/organizations.
// Currently we allow only one team per account.
type Account struct {
	ID          string      `db:"admin_id"`
	Email       string      `db:"email"`
	DisplayName null.String `db:"display_name"`
	CreatedUTC  chrono.Time `db:"created_utc"`
	UpdatedUTC  chrono.Time `db:"updated_utc"`
	Team        null.String `db:"team"` // The id of the team associated with this account.
}
