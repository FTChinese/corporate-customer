package admin

import (
	"encoding/json"
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// BaseLicence is the licence a team purchased.
// A team can purchase multiple licences, as they usually do.
// After purchasing licences, the team can send
// invitations to team member for a specific licence.
// Once a team member accepts the invitation, this licence
// is assigned to this member and cannot be reassigned
// unless it is revoked.
// When a licence is revoked, backup reader's current
// current membership status and delete that row, clear
// the AssigneeID field.
type BaseLicence struct {
	ID         string      `json:"id" db:"licence_id"`
	TeamID     string      `json:"team_id" db:"team_id"`
	AssigneeID null.String `json:"-" db:"assignee_id"` // Only exists after reader accepted an invitation.
	ExpireDate chrono.Date `json:"expireDate" db:"expire_date"`
	Status     null.String `json:"status" db:"current_status"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

// IsAvailable checks whether the licence is
// assigned to someone else.
func (l BaseLicence) IsAvailable() bool {
	return l.AssigneeID.IsZero()
}

func (l *BaseLicence) AssignTo(ftcID string) {
	l.AssigneeID = null.StringFrom(ftcID)
}

// ExpandedLicence includes the data of the licence,
// its attached plan, a latest invitation if present
// and the assignee if it is already granted.
type ExpandedLicence struct {
	BaseLicence
	Plan       plan.BasePlan
	Invitation Invitation // If no one is invited to use this licence, its fields are empty.
	Assignee   Assignee   // If no use is granted to use this licence, its fields are empty.
}

// LicenceList is used to output a list of licence with pagination.
type LicenceList struct {
	Total int64             `json:"total"` // The total number of rows.
	Data  []ExpandedLicence `json:"data"`
}

// LicenceSchema is used to save/retrieve licence.
type LicenceSchema struct {
	BaseLicence
	CurrentPlan    string      `db:"current_plan"`    // The raw JSON column is a string.
	LastInvitation null.String `db:"last_invitation"` // Raw json column.
	Assignee
}

// ExpandedLicence transforms a row retrieved from DB to output format.
func (s LicenceSchema) ExpandedLicence() (ExpandedLicence, error) {
	var p plan.BasePlan
	var inv Invitation

	if err := json.Unmarshal([]byte(s.CurrentPlan), &p); err != nil {
		return ExpandedLicence{}, err
	}

	if s.LastInvitation.Valid {
		if err := json.Unmarshal([]byte(s.LastInvitation.String), &inv); err != nil {
			return ExpandedLicence{}, err
		}
	}

	return ExpandedLicence{
		BaseLicence: s.BaseLicence,
		Plan:        p,
		Invitation:  inv,
		Assignee:    s.Assignee,
	}, nil
}
