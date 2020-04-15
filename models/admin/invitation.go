package admin

import (
	"github.com/FTChinese/go-rest/chrono"
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
	ID            string      `db:"invitation_id"`
	LicenceID     string      `db:"licence_id"`
	TeamID        string      `db:"team_id"`
	Token         string      `db:"token"`
	ExpiresInDays int64       `db:"expires_in_days"`
	Description   null.String `db:"description"`
	Accepted      bool        `db:"accepted"`
	Revoked       bool        `db:"revoked"`
	CreatedUTC    chrono.Time `db:"created_utc"`
	UpdatedUTC    chrono.Time `db:"updated_utc"`
}

func (i *Invitation) Accept() {
	i.Accepted = true
}

func (i Invitation) Expired() bool {
	now := time.Now().Unix()

	created := i.CreatedUTC.Time.Unix()

	// Default 7 days * 24 * 60 * 60
	return (created + i.ExpiresInDays*86400) < now
}

func (i Invitation) IsValid() bool {
	return !i.Expired() && !i.Revoked && !i.Accepted
}

type ExpandedInvitation struct {
	Invitation
	Assignee Assignee
}

type InvitationSchema struct {
	Invitation
	Assignee
}

func (s InvitationSchema) ExpandedInvitation() ExpandedInvitation {
	return ExpandedInvitation{
		Invitation: s.Invitation,
		Assignee:   s.Assignee,
	}
}
