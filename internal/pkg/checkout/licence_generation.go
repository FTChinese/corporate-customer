package checkout

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/ftacademy/pkg/subs"
	"github.com/FTChinese/go-rest/enum"
)

// LicenceGenParams contains the data required to generate
// a new licence.
type LicenceGenParams struct {
	Price     price.Price        // The price selected.
	LicTxn    LicenceTransaction // The transaction to create licence
	CurLic    licence.Licence    // Optional current licence only applicable to queue kind == renew.
	Assignee  licence.Assignee   // Optional user of current licence
	CurMember reader.Membership  // Optional current membership to update
}

type LicenceGenerated struct {
	Transaction            LicenceTransaction // Finalized transaction
	LicenceVersion         licence.Versioned  // Licence after modified. If a licence is renewed and then auto-revoked, it should have two records.
	licence.MemberModified                    // Exists only if a licence was already granted and the corresponding membership is renewed.
}

// IsRenewed checks if a granted licence is renewed successfully together with the corresponding membership.
// If it returned true, we should send a notification email.
func (g LicenceGenerated) IsRenewed() bool {
	return g.Transaction.Kind == enum.OrderKindRenew && g.LicenceVersion.PostChange.IsGranted() && !g.LicenceVersion.PostChange.HintGrantMismatch
}

// GenerateLicence creates a new/renewed licence.
// Corresponding membership will be updated for renewal.
// It will return error if
func GenerateLicence(params LicenceGenParams) (LicenceGenerated, error) {
	// Build licence from queue.
	newLic, err := params.LicTxn.BuildLicence(
		params.CurLic,
		params.Price)

	if err != nil {
		return LicenceGenerated{}, err
	}

	// Mark transaction finished.
	// It shouldn't be touched after this operation.
	finalizedTxn := params.LicTxn.Finalize()
	licVer := newLic.
		Versioned(licence.VersionActionFromOrderKind(params.LicTxn.Kind)).
		WithPriorVersion(params.CurLic)

	// If no one is using this licence.
	if !newLic.IsGranted() {
		return LicenceGenerated{
			Transaction:    finalizedTxn,
			LicenceVersion: licVer,
			MemberModified: licence.MemberModified{}, // No related membership
		}, nil
	}

	// The licence is granted to someone.
	// If there's a mismatch between membership and the licence.
	if !newLic.IsGrantedTo(params.CurMember) {
		return LicenceGenerated{
			Transaction: finalizedTxn,
			LicenceVersion: licVer.
				WithMismatched(params.CurMember),
			// AnteChange is not modified.
			MemberModified: licence.MemberModified{},
		}, nil
	}

	// Now licence is granted to this membership.
	// We have to make sure that its current membership
	// can be renewed.
	// If user manually switched to IAP/Stripe after b2b
	// expired, licence should be reset and keep membership intact.
	// If the granted membership switched to IAP/Stripe,
	// the generated licence should be reset to initial state as if it is revoked.
	// We should also record that this licence is automatically revoked due to renewal inconsistency.
	mm, err := newLic.Grant(params.Assignee, params.CurMember)
	if err != nil {
		if errors.Is(err, subs.ErrOverrideAutoRenewForbidden) {
			return LicenceGenerated{
				Transaction: finalizedTxn,
				LicenceVersion: licVer.
					WithMismatched(params.CurMember),
				MemberModified: licence.MemberModified{},
			}, nil
		}

		return LicenceGenerated{}, err
	}

	return LicenceGenerated{
		Transaction:    finalizedTxn,
		LicenceVersion: licVer.WithMembershipVersioned(mm.MembershipVersion.ID),
		MemberModified: mm,
	}, nil
}
