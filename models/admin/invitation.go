package admin

import (
	"github.com/FTChinese/b2b/models/validator"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"strings"
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
	InviteeEmail  string      `db:"invitee_email"`
	ExpiresInDays int64       `db:"expires_in_days"`
	Description   null.String `db:"description"`
	Accepted      bool        `db:"accepted"`
	Revoked       bool        `db:"revoked"`
	CreatedUTC    chrono.Time `db:"created_utc"`
	UpdatedUTC    chrono.Time `db:"updated_utc"`
}

func NewInvitation(f InvitationForm) (Invitation, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return Invitation{}, err
	}
	return Invitation{
		ID:           "inv_" + rand.String(12),
		LicenceID:    "",
		TeamID:       "",
		Token:        token,
		InviteeEmail: f.Email,
		Description:  null.NewString(f.Description, f.Description != ""),
	}, nil
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

// ExpandedInvitation contains the details of an invitation,
// which licence it is grating, and which team send the
// invitation.
// This is used to show a single invitation to the admin,
// or when we are verifying a user's attempt to accept the
// invitation.
// Using a JOIN to retrieve the invitation and licence by
// licence id, and then retrieve the team and the
// licence's plan, which should be in the cache.
type ExpandedInvitation struct {
	Invitation
	Licence ExpandedLicence
	Team    Team
}

type InvitationForm struct {
	Email       string `form:"email"`
	Description string `form:"description"`
	Errors      map[string]string
}

func (f *InvitationForm) Validate() bool {
	f.Email = strings.TrimSpace(f.Email)
	f.Description = strings.TrimSpace(f.Description)

	msg := validator.New("邮箱").Required().Email().Validate(f.Email)
	if msg != "" {
		f.Errors["email"] = msg
	}

	msg = validator.New("备注").Max(512).Validate(f.Description)
	if msg != "" {
		f.Errors["description"] = msg
	}

	return len(f.Errors) != 0
}
