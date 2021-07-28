package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type GrantResult struct {
	Licence    Licence               `json:"licence"` // The licence after grated.
	Membership reader.Membership     `json:"membership"`
	Snapshot   reader.MemberSnapshot `json:"-"`
	Invoice    reader.Invoice        `json:"invoice"`
}

// NewGrantResult builds data related to granting a licence
// to an assignee.
// The assignee might already have a membership, or might not.
// If the assignee already has a valid membership, it could
// only be allowed to come from alipay or wechat.
// IAP and stripe should not be allowed to use a licence.
// @param grantedLic - the licence to grant
// @param m - current membership
func NewGrantResult(grantedLic Licence, m reader.Membership) GrantResult {
	var inv reader.Invoice
	// If membership is still valid, turn remaining days to
	// carried-over addon.
	if !m.IsExpired() {
		inv = m.CarryOverInvoice()
	}

	newM := NewMembership(
		reader.UserIDs{
			CompoundID: grantedLic.Assignee.FtcID.String,
			FtcID:      grantedLic.Assignee.FtcID,
			UnionID:    grantedLic.Assignee.UnionID,
		},
		grantedLic.BaseLicence,
		m.AddOn.Plus(addon.New(inv.Tier, inv.TotalDays())),
	)

	var snapshot reader.MemberSnapshot
	if !m.IsZero() {
		snapshot = m.Archive(reader.B2BArchiver(reader.ArchiveActionGrant))
	}

	return GrantResult{
		Licence:    grantedLic,
		Membership: newM,
		Snapshot:   snapshot,
		Invoice:    inv,
	}
}

type RevokeResult struct {
	Licence    Licence               `json:"licence"`
	Membership reader.Membership     `json:"membership"`
	Snapshot   reader.MemberSnapshot `json:"snapshot"`
}

// RevokeLicence revokes a licence granted to a membership.
// Ideally, if membership has addon, we should use addon
// invoices to re-build membership.
// The process, however, is quite complicated and we have
// duplicate all invoice manipulate from subscription api.
// To save effort, we simply change the expiration date
// to now and set payment method to a one-time payment,
// and leave addon fields untouched so that when client
// detects addon should be transferred, api to handle it.
func RevokeLicence(m reader.Membership) reader.Membership {
	var tier enum.Tier
	if m.HasAddOn() {
		if m.AddOn.Standard > 0 {
			tier = enum.TierStandard
		} else if m.AddOn.Premium > 0 {
			tier = enum.TierPremium
		}
	} else {
		tier = m.Tier
	}

	return reader.Membership{
		UserIDs: m.UserIDs,
		Edition: price.Edition{
			Tier:  tier,
			Cycle: m.Cycle,
		},
		LegacyTier:    null.Int{},
		LegacyExpire:  null.Int{},
		ExpireDate:    chrono.DateNow(),
		PaymentMethod: enum.PayMethodAli,
		FtcPlanID:     null.String{},
		StripeSubsID:  null.String{},
		StripePlanID:  null.String{},
		AutoRenewal:   false,
		Status:        enum.SubsStatusNull,
		AppleSubsID:   null.String{},
		B2BLicenceID:  null.String{},
		AddOn:         m.AddOn,
	}.Sync()
}
