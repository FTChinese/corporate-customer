package invitee

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"time"
)

func GenerateMemberID() string {
	return "mmb_" + rand.String(12)
}

type Membership struct {
	SubsID         null.String     `db:"subs_id"`
	SubsCompoundID null.String     `db:"subs_compound_id"`
	SubsFtcID      null.String     `db:"subs_ftc_id"`
	SubsUnionID    null.String     `db:"subs_union_id"`
	LegacyWxID     null.String     `db:"legacy_wx_id"`
	LegacyTier     null.Int        `db:"legacy_tier"` // 10 - standard, 100 - premium
	LegacyExpire   null.Int        `db:"legacy_expire"`
	Tier           enum.Tier       `db:"tier"`
	Cycle          enum.Cycle      `db:"cycle"`
	ExpireDate     chrono.Date     `db:"expire_date"`
	AutoRenew      bool            `db:"auto_renew"`
	PayMethod      enum.PayMethod  `db:"payment_method"`
	StripSubsID    null.String     `db:"stripe_subs_id"`
	StripePlanID   null.String     `db:"stripe_plan_id"`
	SubsStatus     enum.SubsStatus `db:"subs_status"`
	AppleSubsID    null.String     `db:"apple_subs_id"`
}

func (m Membership) BuildOn(l admin.ExpandedLicence) Membership {
	if m.IsZero() {
		m.SubsID = null.StringFrom(GenerateMemberID())
		m.SubsCompoundID = l.AssigneeID
		m.SubsFtcID = l.AssigneeID
	}

	if m.SubsID.IsZero() {
		m.SubsID = null.StringFrom(GenerateMemberID())
	}

	// Erase the legacy fields so that
	// we could call Normalize correctly.
	m.LegacyTier = null.Int{}
	m.LegacyExpire = null.Int{}
	m.Tier = l.Plan.Tier
	m.Cycle = l.Plan.Cycle
	m.ExpireDate = l.ExpireDate
	m.AutoRenew = false
	m.PayMethod = enum.PayMethodB2B
	m.StripSubsID = null.String{}
	m.StripePlanID = null.String{}
	m.SubsStatus = enum.SubsStatusNull
	m.AppleSubsID = null.String{}

	m.Normalize()

	return m
}

func (m Membership) IsZero() bool {
	return m.SubsCompoundID.IsZero()
}

func (m *Membership) Normalize() {
	if m.IsZero() {
		return
	}

	// Indicating using the legacy id columns
	if m.SubsFtcID.IsZero() && m.SubsUnionID.IsZero() {
		// Indicates a wechat-only membership
		if m.SubsCompoundID.String == m.LegacyWxID.String {
			m.SubsUnionID = m.LegacyWxID
		} else { // Indicates a ftc membership
			m.SubsFtcID = m.SubsCompoundID
		}
	}

	// If using the legacy expiration column only.
	if m.LegacyExpire.Valid && m.ExpireDate.IsZero() {
		m.ExpireDate = chrono.DateFrom(time.Unix(m.LegacyExpire.Int64, 0))
	}

	// If using the new expiration column only.
	if !m.ExpireDate.IsZero() && m.LegacyExpire.IsZero() {
		m.LegacyExpire = null.IntFrom(m.ExpireDate.Unix())
	}

	if m.LegacyTier.Valid && m.Tier == enum.TierNull {
		switch m.LegacyTier.Int64 {
		case 10:
			m.Tier = enum.TierStandard
		case 100:
			m.Tier = enum.TierPremium
		}
	}

	if m.Tier != enum.TierNull && m.LegacyTier.IsZero() {
		switch m.Tier {
		case enum.TierStandard:
			m.LegacyTier = null.IntFrom(10)
		case enum.TierPremium:
			m.LegacyTier = null.IntFrom(100)

		}
	}
}

func (m Membership) IsExpired() bool {
	if m.IsZero() {
		return true
	}

	// Ignore expire date for auto renewal subscriptions.
	if m.AutoRenew {
		return false
	}

	return m.ExpireDate.Before(time.Now().Truncate(24 * time.Hour))
}

func GenerateSnapshotID() string {
	return "snp_" + rand.String(12)
}

type MemberSnapshot struct {
	SnapshotID string              `db:"snapshot_id"`
	Reason     enum.SnapshotReason `db:"reason"`
	CreatedUTC chrono.Time         `db:"created_utc"`
	Membership
}
