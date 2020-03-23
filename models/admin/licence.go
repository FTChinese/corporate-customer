package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Licence is the subscription a team purchased.
// A team can purchase multiple licences, as they often do.
// After purchasing licences, the team can send
// invitations to team member for a specific licence.
// Once a team member accepts the invitation, this licence
// is assigned to this member and cannot be reassigned
// unless it is revoked.
// When a licence is revoked, backup reader's current
// current membership status and delete that row, clear
// the AssigneeID field.
type Licence struct {
	ID         string      `db:"licence_id"`
	Tier       enum.Tier   `db:"tier"`
	Cycle      enum.Cycle  `db:"cycle"`
	ExpireDate chrono.Date `db:"expire_date"`
	AssigneeID null.String `db:"assignee_id"`
	CreatedUTC chrono.Time `db:"created_utc"`
	UpdatedUTC chrono.Time `db:"updated_utc"`
	Team       string      `db:"team"`
}
