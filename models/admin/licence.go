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
	ID         string        `json:"id" db:"licence_id"`
	TeamID     string        `json:"team_id" db:"team_id"`
	AssigneeID null.String   `json:"-" db:"assignee_id"` // Only exists after reader accepted an invitation.
	ExpireDate chrono.Date   `json:"expireDate" db:"expire_date"`
	Status     LicenceStatus `json:"status" db:"current_status"`
	CreatedUTC chrono.Time   `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time   `json:"updatedUtc" db:"updated_utc"`
}

// IsAvailable checks whether the licence is
// assigned to someone else.
func (l BaseLicence) IsAvailable() bool {
	return l.AssigneeID.IsZero() && l.Status == LicStatusAvailable
}

func (l *BaseLicence) AssignTo(ftcID string) {
	l.AssigneeID = null.StringFrom(ftcID)
}

// Licence contains all data of a licence.
type Licence struct {
	BaseLicence
	Plan       plan.BasePlan
	Invitation Invitation // If no one is invited to use this licence, its fields are empty.
}

func (l Licence) WithInvitation(inv Invitation) Licence {
	l.Status = LicStatusInvited
	l.Invitation = inv
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// LicenceSchema defines the DB schema for licence table.
type LicenceSchema struct {
	BaseLicence
	CurrentPlan    string      `db:"current_plan"`    // The raw JSON column is a string.
	LastInvitation null.String `db:"last_invitation"` // Raw json column.
}

// Licence create the Licence instance from raw db data.
func (ls LicenceSchema) Licence() (Licence, error) {
	var p plan.BasePlan
	var inv Invitation

	if err := json.Unmarshal([]byte(ls.CurrentPlan), &p); err != nil {
		return Licence{}, err
	}

	if ls.LastInvitation.Valid {
		if err := json.Unmarshal([]byte(ls.LastInvitation.String), &inv); err != nil {
			return Licence{}, err
		}
	}

	return Licence{
		BaseLicence: ls.BaseLicence,
		Plan:        p,
		Invitation:  inv,
	}, nil
}

func (ls LicenceSchema) WithInvitationSent(inv Invitation) (LicenceSchema, error) {
	invDoc, err := json.Marshal(inv)
	if err != nil {
		return ls, err
	}

	ls.Status = LicStatusInvited
	ls.LastInvitation = null.StringFrom(string(invDoc))
	ls.UpdatedUTC = chrono.TimeNow()

	return ls, nil
}

// WithInvitationRevoked revokes an invitation of a licence
// so that admin could invite another one to use this licence.
// This is used after invitation is sent but before it is
// accepted.
func (ls LicenceSchema) WithInvitationRevoked() LicenceSchema {
	ls.Status = LicStatusAvailable
	ls.LastInvitation = null.String{}
	ls.UpdatedUTC = chrono.TimeNow()

	return ls
}

// Revoke unlink a user from a licence.
func (ls LicenceSchema) Revoke() LicenceSchema {
	ls.AssigneeID = null.String{}
	ls.Status = LicStatusAvailable
	ls.UpdatedUTC = chrono.TimeNow()

	return ls
}

// ExpandedLicence includes the data of the licence,
// its attached plan, a latest invitation if present
// and the assignee if it is already granted.
type ExpandedLicence struct {
	Licence
	Assignee Assignee // If no use is granted to use this licence, its fields are empty.
}

// ExpLicenceSchema is used to save/retrieve licence.
type ExpLicenceSchema struct {
	LicenceSchema
	Assignee
}

// ExpandedLicence transforms a row retrieved from DB to output format.
func (s ExpLicenceSchema) ExpandedLicence() (ExpandedLicence, error) {

	lic, err := s.Licence()
	if err != nil {
		return ExpandedLicence{}, err
	}

	return ExpandedLicence{
		Licence:  lic,
		Assignee: s.Assignee,
	}, nil
}

// ExpLicenceList is used to output a list of licence with pagination.
type ExpLicenceList struct {
	Total int64             `json:"total"` // The total number of rows.
	Data  []ExpandedLicence `json:"data"`
}
