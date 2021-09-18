package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/subs"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// MemberModified contains the data related to modification of membership.
type MemberModified struct {
	MembershipVersion  reader.MembershipVersioned `json:"membershipVersion"`  // Membership prior to modification.
	CarryOverInvoice   reader.Invoice             `json:"carryOverInvoice"`   // Optional addon invoice in case previous remaining period is converted to addon.
	IsSwitchingLicence bool                       `json:"isSwitchingLicence"` // Admin might decide to change licence granted to a member. In such case we should revoke previous licence after granting a new one.
}

// NewMembership builds a new membership based on a licence.
// You cannot use licence's assignee id here for 2 reasons:
// * The licence might not be granted to anyone, thus it's empty;
// * User might have wechat id linked, while  the assignee id only reflects email account.
func (l Licence) NewMembership(userIDs reader.UserIDs, addOn addon.AddOn) reader.Membership {
	return reader.Membership{
		UserIDs:       userIDs,
		Edition:       l.Edition,
		LegacyTier:    null.Int{},
		LegacyExpire:  null.Int{},
		ExpireDate:    chrono.DateFrom(l.CurrentPeriodEndUTC.Time),
		PaymentMethod: enum.PayMethodB2B,
		FtcPlanID:     null.StringFrom(l.LatestPrice.ID),
		StripeSubsID:  null.String{},
		StripePlanID:  null.String{},
		AutoRenewal:   false,
		Status:        enum.SubsStatusNull,
		AppleSubsID:   null.String{},
		B2BLicenceID:  null.StringFrom(l.ID),
		AddOn:         addOn,
		VIP:           false,
	}.Sync()
}

// Grant modifies membership based on a latest licence..
// The assignee might already have a membership, or might not.
// If the assignee already has a valid membership, it could
// only be allowed to come from alipay or wechat.
// IAP and stripe should not be allowed to use a licence.
// @param to - to whom the licence is granted.
// @param curMmb - current membership used to calculate addon
func (l Licence) Grant(to Assignee, curMmb reader.Membership) (MemberModified, error) {
	subsKind, err := l.SubsKind(curMmb)
	if err != nil {
		return MemberModified{}, err
	}

	// Optional carry-over invoice.
	var inv reader.Invoice
	if subsKind == subs.KindOverrideOneTime {
		inv = curMmb.CarryOverInvoice().
			WithLicTxID(l.LatestTransactionID)
	}

	// Build new membership.
	memberVersioned := l.NewMembership(
		reader.UserIDs{
			CompoundID: to.FtcID.String,
			FtcID:      to.FtcID,
			UnionID:    to.UnionID,
		},
		curMmb.AddOn.Plus(addon.New(inv.Tier, inv.TotalDays())),
	).
		Version(reader.B2BArchiver(reader.ArchiveActionGrant)).
		WithPriorVersion(curMmb).
		WithB2BTxnID(l.LatestTransactionID.String)

	return MemberModified{
		MembershipVersion:  memberVersioned,
		CarryOverInvoice:   inv,
		IsSwitchingLicence: subsKind == subs.KindB2BSwitchLicence,
	}, nil
}
