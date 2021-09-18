package reader

import (
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"math"
	"time"
)

var tierToCode = map[enum.Tier]int64{
	enum.TierStandard: 10,
	enum.TierPremium:  100,
}

var codeToTier = map[int64]enum.Tier{
	10:  enum.TierStandard,
	100: enum.TierPremium,
}

func GetTierCode(tier enum.Tier) int64 {
	return tierToCode[tier]
}

// Membership contains a user's membership details
// This is actually called subscription by Stripe.
// A membership might be create from various sources:
// * Alipay / Wxpay - Classified under FTC retail
// * B2B
// * Stripe
// * Apple IAP
// We should keep those sources mutually exclusive.
type Membership struct {
	UserIDs
	price.Edition
	LegacyTier    null.Int       `json:"-" db:"vip_type"`
	LegacyExpire  null.Int       `json:"-" db:"expire_time"`
	ExpireDate    chrono.Date    `json:"expireDate" db:"expire_date"`
	PaymentMethod enum.PayMethod `json:"payMethod" db:"payment_method"`
	FtcPlanID     null.String    `json:"ftcPlanId" db:"ftc_plan_id"`
	StripeSubsID  null.String    `json:"stripeSubsId" db:"stripe_subs_id"`
	StripePlanID  null.String    `json:"stripePlanId" db:"stripe_plan_id"`
	AutoRenewal   bool           `json:"autoRenew" db:"auto_renewal"`
	// This is used to save stripe subscription status.
	// Since wechat and alipay treats everything as one-time purchase, they do not have a complex state machine.
	// If we could integrate apple in-app purchase, this column
	// might be extended to apple users.
	// Only `active` should be treated as valid member.
	// Wechat and alipay defaults to `active` for backward compatibility.
	Status       enum.SubsStatus `json:"status" db:"subs_status"`
	AppleSubsID  null.String     `json:"appleSubsId" db:"apple_subs_id"`
	B2BLicenceID null.String     `json:"b2bLicenceId" db:"b2b_licence_id"`
	addon.AddOn
	VIP bool `json:"vip" db:"is_vip"`
}

// IsZero test whether the instance is empty.
func (m Membership) IsZero() bool {
	return m.CompoundID == "" && m.Tier == enum.TierNull
}

// IsExpired tests if the membership's expiration date is before now.
// A non-existing membership is treated as expired.
// Auto renewal is treated as not expired.
func (m Membership) IsExpired() bool {
	// If membership does not exist, it is treated as expired.
	if m.IsZero() {
		return true
	}

	// If expire date is before now, AND auto renew is false,
	// we treat this one as actually expired.
	// If ExpireDate is passed, but auto renew is true, we still
	// treat this one as not expired.
	return m.ExpireDate.Before(time.Now().Truncate(24*time.Hour)) && !m.AutoRenewal
}

// isLegacyOnly checks whether the edition information only comes from
// LegacyTier and LegacyExpire fields.
func (m Membership) isLegacyOnly() bool {
	if m.LegacyExpire.Valid && m.LegacyTier.Valid && m.ExpireDate.IsZero() && m.Tier == enum.TierNull {
		return true
	}

	return false
}

// isAPIOnly checks whether the edition information only comes from
// Tier and Cycle fields.
func (m Membership) isAPIOnly() bool {
	if (m.LegacyExpire.IsZero() && m.LegacyTier.IsZero()) && (!m.ExpireDate.IsZero() && m.Tier != enum.TierNull) {
		return true
	}

	return false
}

// RemainingDays calculates how many day left up until now.
// If the returned days is less than 0, the membership is expired
// if it is not auto renewable.
func (m Membership) RemainingDays() int64 {
	h := time.Until(m.ExpireDate.Time).Hours()

	return int64(math.Ceil(h / 24))
}

// Sync turns legacy vip_type and expire_time into
// member_tier and expire_date columns, or vice versus.
func (m Membership) Sync() Membership {
	if m.IsZero() {
		return m
	}

	// Syn from legacy format to api created columns
	if m.isLegacyOnly() {
		// Note the conversion is not exactly the same moment since Golang converts Unix in local time.
		expireDate := time.Unix(m.LegacyExpire.Int64, 0)

		m.ExpireDate = chrono.DateFrom(expireDate)
		m.Tier = codeToTier[m.LegacyTier.Int64]
		// m.Cycle cannot be determined

		return m
	}

	// Sync from api columns to legacy column
	if m.isAPIOnly() {
		m.LegacyExpire = null.IntFrom(m.ExpireDate.Unix())
		m.LegacyTier = null.IntFrom(tierToCode[m.Tier])

		return m
	}

	// Otherwise do not touch it.
	// PROBLEM: if both sides have data but they are not synced, the discrepancy will be carried forward.
	return m
}

// IsOneTime checks whether current membership is purchased
// via alipay or wechat pay.
func (m Membership) IsOneTime() bool {
	// For backward compatibility. If Tier field comes from LegacyTier, then PayMethod field will be null.
	// We treat all those cases as wxpay or alipay.
	if m.Tier != enum.TierNull && m.PaymentMethod == enum.PayMethodNull {
		return true
	}

	return m.PaymentMethod == enum.PayMethodAli || m.PaymentMethod == enum.PayMethodWx
}

func (m Membership) IsAutoRenew() bool {
	return m.PaymentMethod == enum.PayMethodApple || m.PaymentMethod == enum.PayMethodStripe
}

func (m Membership) IsIAP() bool {
	return !m.IsZero() && m.PaymentMethod == enum.PayMethodApple && m.AppleSubsID.Valid
}

func (m Membership) IsStripe() bool {
	return !m.IsZero() && m.PaymentMethod == enum.PayMethodStripe && m.StripeSubsID.Valid
}

func (m Membership) IsB2B() bool {
	return m.PaymentMethod == enum.PayMethodB2B && m.B2BLicenceID.Valid
}

// LicenceRevoked change membership to expired state.
func (m Membership) LicenceRevoked() Membership {
	if m.HasAddOn() {
		if m.AddOn.Standard > 0 {
			m.Tier = enum.TierStandard
		} else if m.AddOn.Premium > 0 {
			m.Tier = enum.TierPremium
		}
	}

	m.ExpireDate = chrono.DateFrom(time.Now().AddDate(0, 0, -1))
	m.PaymentMethod = enum.PayMethodAli
	return m
}
