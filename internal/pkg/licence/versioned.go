package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type VersionAction string

const (
	VersionActionNull   VersionAction = ""
	VersionActionCreate VersionAction = "create"
	VersionActionRenew  VersionAction = "renew"
	VersionActionGrant  VersionAction = "grant"
	VersionActionRevoke VersionAction = "revoke"
)

func VersionActionFromOrderKind(k enum.OrderKind) VersionAction {
	switch k {
	case enum.OrderKindCreate:
		return VersionActionCreate

	case enum.OrderKindRenew:
		return VersionActionRenew

	default:
		return VersionActionNull
	}
}

// Versioned takes a snapshot after a licence is created/renewed/granted/revoked.
// When a granted licence is renewed, it is possible that
// the assignee already changed to other payment methods and
// the granting should be revoked.
type Versioned struct {
	Action              VersionAction         `json:"action" db:"action_kind"` // Action type
	AnteChange          LicJSON               `json:"anteChange" db:"ante_change"`
	MembershipVersionID null.String           `json:"membershipVersionId" db:"membership_version_id"`
	MismatchedMember    reader.MembershipJSON `json:"mismatchedMember" db:"mismatched_member"` // If a licence is granted to a user but user has changed subscription with a higher priority, mark this licence as mismatched so that admin could revoke it directly.
	PostChange          LicJSON               `json:"postChange" db:"post_change"`             // Licence after modification
	CreatedUTC          chrono.Time           `json:"createdUtc" db:"created_utc"`
}

// WithPriorVersion saves licence prior to being changed.
// The prior version might be empty is the licence is newly created.
func (v Versioned) WithPriorVersion(l Licence) Versioned {
	if l.IsZero() {
		return v
	}

	v.AnteChange = LicJSON{l}
	return v
}

// WithMembershipVersioned keeps id of membership snapshot.
// This happens when membership is changed by licence.
// The passed-in id is the ID field of reader.MembershipVersioned
func (v Versioned) WithMembershipVersioned(id string) Versioned {
	v.MembershipVersionID = null.StringFrom(id)
	return v
}

// WithMismatched saves a snapshot of membership when a licence
// is found to be granted to a membership with higher priority.
// This tells why the licence is auto-revoked.
func (v Versioned) WithMismatched(m reader.Membership) Versioned {
	v.PostChange = LicJSON{v.PostChange.WithGrantMismatch(true)}
	v.MismatchedMember = reader.MembershipJSON{Membership: m}
	return v
}

// Versioned takes a snapshot after a licence changed.
// If the revoking is performed automatically, caused by a mismatch between the licence and membership,
// pass current membership to record why the membership
// is causing the revoke; otherwise pass empty membership.
func (l Licence) Versioned(k VersionAction) Versioned {
	return Versioned{
		Action:              k,
		AnteChange:          LicJSON{},               // Empty for newly created licence
		MembershipVersionID: null.String{},           // Empty is no linked membership touched.
		MismatchedMember:    reader.MembershipJSON{}, // Only exists when licence is granted but assignee is not using it.
		PostChange:          LicJSON{Licence: l},
		CreatedUTC:          chrono.TimeNow(),
	}
}
