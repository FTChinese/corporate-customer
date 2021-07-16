package licence

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

func NewMembership(userIDs reader.UserIDs, l BaseLicence, addOn addon.AddOn) reader.Membership {
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
	}
}
