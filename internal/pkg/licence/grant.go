package licence

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/subs"
	"github.com/FTChinese/go-rest/chrono"
)

// SubsKind deduces what kind of action is taking when
// granting a licence to an optional existing membership.
// Returns:
// * subs.KindB2BNew if current membership is empty or expired;
// * subs.KindOverrideOneTime if current membership is valid and comes from one-time purchase
// * subs.KindB2BRenew if current membership comes from this licence
// * subs.KindB2BSwitchLicence if current membership comes from B2B but the licence-to-grant differs from the one in use.
// Return subs.ErrOverrideAutoRenewForbidden if current membership is an auto-renewal one.
func (l Licence) SubsKind(m reader.Membership) (subs.Kind, error) {
	// Expired membership is treated as a new one.
	if m.IsExpired() {
		return subs.KindB2BNew, nil
	}

	// You cannot grant b2b to auto-renewal user.
	if m.IsAutoRenew() {
		return subs.KindZero, subs.ErrOverrideAutoRenewForbidden
	}

	// One-time purchase will be overridden.
	if m.IsOneTime() {
		return subs.KindOverrideOneTime, nil
	}

	// Renewing existing one, or change licence.
	if m.IsB2B() {
		if l.ID == m.B2BLicenceID.String {
			return subs.KindB2BRenew, nil
		}

		// When admin is granting a new licence to an existing
		// b2b user, we should revoke old licence.
		return subs.KindB2BSwitchLicence, nil
	}

	return subs.KindZero, errors.New("licence granting failed due to unknown current membership status of this user")
}

// WithGranted changes a licence status to granted with assignee
// field populated and invitation field updated.
func (l Licence) WithGranted(to Assignee, inv Invitation) Licence {
	l.HintGrantMismatch = false
	l.Status = LicStatusGranted
	l.LatestInvitation = InvitationJSON{inv.Accepted()}
	l.AssigneeID = to.FtcID
	l.UpdatedUTC = chrono.TimeUTCNow()

	return l
}

// IsGrantedTo checks if a licence is granted to the
// specified membership.
func (l Licence) IsGrantedTo(m reader.Membership) bool {
	if !m.IsB2B() {
		return false
	}

	return l.ID == m.B2BLicenceID.String
}

// IsGranted tests whether a licence is granted to any reader.
func (l Licence) IsGranted() bool {
	return l.Status == LicStatusGranted && l.AssigneeID.Valid
}

// WithGrantMismatch set HintGrantMismatch to true as an indicator
// to tell admin possible problems with the licence.
func (l Licence) WithGrantMismatch(mismatched bool) Licence {
	l.HintGrantMismatch = mismatched
	return l
}

// GrantResult contains all data generated after grating
// a licence to a user.
type GrantResult struct {
	LicenceVersion Versioned `json:"licenceVersion"`
	MemberModified
}

type GrantParams struct {
	CurLic    Licence           // Licence to be granted
	CurInv    Invitation        // Invitation user is accepting
	To        Assignee          // To whom the licence the granted
	CurMember reader.Membership // Current membership
}

// GrantLicence grants a licence to an assignee.
// Return a GrantResult instance or subs.ErrOverrideAutoRenewForbidden.
func GrantLicence(params GrantParams) (GrantResult, error) {
	grantedLic := params.CurLic.
		WithGranted(params.To, params.CurInv.Accepted())

	mm, err := grantedLic.Grant(params.To, params.CurMember)
	if err != nil {
		return GrantResult{}, err
	}

	return GrantResult{
		LicenceVersion: grantedLic.
			Versioned(VersionActionGrant).
			WithPriorVersion(params.CurLic).
			WithMembershipVersioned(mm.MembershipVersion.ID),
		MemberModified: mm,
	}, nil
}
