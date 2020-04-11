package admin

import (
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// Licence is the licence a team purchased.
// Any licence must have a plan attached to it,
// and it must belong to a term.
// Team -> Licence and Plan -> Licence are both
// one-to-many relations.
// A team can purchase multiple licences, as they usually do.
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
	TeamID     string      `db:"team_id"`
	AssigneeID null.String `db:"assignee_id"` // Only exists after reader accepted an invitation.
	ExpireDate chrono.Date `db:"expire_date"`
	Active     bool        `db:"is_active"` // Only active after payment received. This is not controlled by the admin.
	CreatedUTC chrono.Time `db:"created_utc"`
	UpdatedUTC chrono.Time `db:"updated_utc"`
}

// IsAvailable checks whether the licence is
// assigned to someone else.
func (l Licence) IsAvailable() bool {
	return l.AssigneeID.IsZero() && l.Active
}

func (l *Licence) AssignTo(ftcID string) {
	l.AssigneeID = null.StringFrom(ftcID)
}

// ExpandedLicence includes the plan of this licence
// and optional assignee.
type ExpandedLicence struct {
	Licence
	Plan     plan.BasePlan
	Assignee Assignee
}

type LicenceSchema struct {
	Licence
	plan.BasePlan
	Assignee
}

func (s LicenceSchema) ExpandedLicence() ExpandedLicence {
	return ExpandedLicence{
		Licence:  s.Licence,
		Plan:     s.BasePlan,
		Assignee: s.Assignee,
	}
}
