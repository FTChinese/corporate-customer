package model

import (
	"encoding/json"
	"github.com/FTChinese/ftacademy/internal/pkg/plan"
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
	ID               string        `json:"id" db:"licence_id"`
	TeamID           string        `json:"teamId" db:"team_id"`
	AssigneeID       null.String   `json:"-" db:"assignee_id"` // Only exists after reader accepted an invitation.
	ExpireDate       chrono.Date   `json:"expireDate" db:"expire_date"`
	TrialStart       chrono.Date   `json:"trialStart" db:"trial_start_date"`
	TrialEnd         chrono.Date   `json:"trialEnd" db:"trial_end_date"`
	Status           LicenceStatus `json:"status" db:"current_status"`
	CreatedUTC       chrono.Time   `json:"createdUtc" db:"created_utc"`
	UpdatedUTC       chrono.Time   `json:"updatedUtc" db:"updated_utc"`
	LastInvitationID null.String   `json:"lastInvitationId" db:"last_invitation_id"`
	LastInviteeEmail null.String   `json:"lastInviteeEmail" db:"last_invitee_email"`
}

// IsAvailable checks whether the licence is
// assigned to someone else.
func (l BaseLicence) IsAvailable() bool {
	return l.Status == LicStatusAvailable && l.AssigneeID.IsZero()
}

// Invited updates invitation status to invited after an invitation is sent.
func (l BaseLicence) Invited(inv Invitation) BaseLicence {
	l.Status = LicStatusInvited
	l.LastInvitationID = null.StringFrom(inv.ID)
	l.LastInviteeEmail = null.StringFrom(inv.Email)
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// CanInvitationBeRevoked ensures that the licence could
// have its invitation revoked.
// A licence could only have its invitation revoked when
// an invitation is sent but not accepted.
func (l BaseLicence) CanInvitationBeRevoked() bool {
	return l.Status == LicStatusInvited && l.AssigneeID.IsZero()
}

// InvitationRevoked revokes an invitation of a licence
// so that admin could invite another one to use this licence.
// This is used after invitation is sent but before it is
// accepted.
func (l BaseLicence) InvitationRevoked() BaseLicence {
	l.Status = LicStatusAvailable
	l.LastInvitationID = null.String{}
	l.LastInviteeEmail = null.String{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

// Revoke unlink a user from a licence.
func (l BaseLicence) Revoke() BaseLicence {
	l.AssigneeID = null.String{}
	l.Status = LicStatusAvailable
	l.LastInvitationID = null.String{}
	l.LastInviteeEmail = null.String{}
	l.UpdatedUTC = chrono.TimeNow()

	return l
}

func (l BaseLicence) CanBeGranted() bool {
	return l.Status == LicStatusInvited && l.AssigneeID.IsZero()
}

// GrantTo links a licence to a ftc reader.
func (l BaseLicence) GrantTo(ftcID string) BaseLicence {
	l.AssigneeID = null.StringFrom(ftcID)
	l.Status = LicStatusGranted
	l.UpdatedUTC = chrono.TimeNow()
	return l
}

// Licence contains all data of a licence.
type Licence struct {
	BaseLicence
	Plan plan.BasePlan
}

// LicenceSchema defines the DB schema for licence table.
type LicenceSchema struct {
	BaseLicence
	CurrentPlan string `json:"-" db:"current_plan"` // The raw JSON column is a string.
}

// Licence create the Licence instance from raw db data.
func (ls LicenceSchema) Licence() (Licence, error) {
	var p plan.BasePlan

	if err := json.Unmarshal([]byte(ls.CurrentPlan), &p); err != nil {
		return Licence{}, err
	}

	return Licence{
		BaseLicence: ls.BaseLicence,
		Plan:        p,
	}, nil
}

// ExpandedLicence includes the data of the licence,
// its attached plan, a latest invitation if present
// and the assignee if it is already granted.
type ExpandedLicence struct {
	Licence
	Assignee Assignee `json:"assignee"` // If no use is granted to use this licence, its fields are empty.
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
type PagedExpLicences struct {
	Total int64             `json:"total"` // The total number of rows.
	Data  []ExpandedLicence `json:"data"`
	Err   error             `json:"-"`
}
