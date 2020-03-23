package admin

import "github.com/FTChinese/go-rest/chrono"

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
	Token         string      `db:"token"`
	InviteeEmail  string      `db:"invitee_email"`
	ExpiresInDays int64       `db:"expires_in_days"`
	Accepted      bool        `db:"accepted"`
	CreatedUTC    chrono.Time `db:"created_utc"`
	UpdatedUTC    chrono.Time `db:"updated_utc"`
	LicenceID     string      `db:"licence_id"`
}
