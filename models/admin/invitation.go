package admin

import (
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"time"
)

// Invitation is an email sent to team member
// to accept a licence.
// An invitee can be sent multiple invitations but
// only one record is kept.
// When user clicked link in the email, we should have enough
// information about this invitation:
// * whether this email already exists, then go to login or signup;
// * whether this user already have a valid subscription;
// * make sure the licence is still valid the moment to assign
// the licence to this email.
//
// If everything are valid, then created/update membership
// for this email, add the AssigneeID field of licence, flag
// this invitation as ready used.
// To achieve this, we need to perform it in one transaction:
// * Flag invitation as used;
// * Add AssigneeID field;
// * Create/Update membership.
//
// Each record is immutable and create a new one
// every time an invitation is sent.
// WHen retrieving the token for verification,
// retrieve only the last created record for this specific
// email.
//
// Upon sending an invitation, we should make sure:
//
// * The licence is not used by anyone else;
// * The InviteeEmail does not have valid subscription.
type Invitation struct {
	ID             string      `json:"id" db:"invitation_id"`
	LicenceID      string      `json:"licenceId" db:"licence_id"`
	TeamID         string      `json:"teamId" db:"team_id"`
	Token          string      `json:"-" db:"token"` // This field is used only when inserting data. Retrieval does not include this field. However, it is included when saving to the JSON column in licence.
	Email          string      `json:"email" db:"email"`
	Description    null.String `json:"description" db:"description"`
	ExpirationDays int64       `json:"expiresInDays" db:"expiration_days"`

	Accepted   bool        `json:"accepted" db:"accepted"`
	Revoked    bool        `json:"revoked" db:"revoked"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

func (i Invitation) Expired() bool {
	now := time.Now().Unix()

	created := i.CreatedUTC.Time.Unix()

	// Default 7 days * 24 * 60 * 60
	return (created + i.ExpirationDays*86400) < now
}

func (i Invitation) IsValid() bool {
	return !i.Expired() && !i.Revoked && !i.Accepted
}

func (i Invitation) Revoke() Invitation {
	i.Revoked = true
	i.UpdatedUTC = chrono.TimeNow()

	return i
}

// InvitationInput contains the essential data client
// submitted to create a new invitation.
type InvitationInput struct {
	Email       string      `json:"email"` // To whom the invitation should be sent.
	Description null.String `json:"description"`
	LicenceID   string      `json:"licenceId"` // Which licence is being granted.
	TeamID      string      `json:"-"`
}

// NewInvitation creates a new Invitation instance based
// on user input.
func (i InvitationInput) NewInvitation() (Invitation, error) {
	token, err := GenerateToken()
	if err != nil {
		return Invitation{}, err
	}

	return Invitation{
		ID:             "inv_" + rand.String(12),
		LicenceID:      i.LicenceID,
		TeamID:         i.TeamID,
		Token:          token,
		Email:          i.Email,
		Description:    i.Description,
		ExpirationDays: 7,
		Accepted:       false,
		Revoked:        false,
		CreatedUTC:     chrono.TimeNow(),
		UpdatedUTC:     chrono.TimeNow(),
	}, nil
}

// InvitedLicence wraps all related information after
// an invitation is created.
type InvitedLicence struct {
	Invitation Invitation
	Licence    BaseLicence   // The licence to grant
	Plan       plan.BasePlan // The plan of this licence
	Assignee   Assignee      // Who will be granted the licence.
}

// InvitationList is used for restful output.
type InvitationList struct {
	Total int64 `json:"total"`
	Data  []Invitation
}
