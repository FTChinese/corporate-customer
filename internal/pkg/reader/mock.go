// +build !production

package reader

import (
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

// MockMembership generates a mocking membership.
// Deprecated
func MockMembership(ftcID string) Membership {
	if ftcID == "" {
		ftcID = uuid.New().String()
	}

	return Membership{
		UserIDs: UserIDs{
			CompoundID: ftcID,
			FtcID:      null.StringFrom(ftcID),
		},
		Edition: price.Edition{
			Tier:  enum.TierStandard,
			Cycle: enum.CycleYear,
		},
		LegacyTier:    null.Int{},
		LegacyExpire:  null.Int{},
		ExpireDate:    chrono.DateUTCFrom(time.Now().AddDate(0, 1, 0)),
		PaymentMethod: enum.PayMethodAli,
		FtcPlanID:     null.String{},
		StripeSubsID:  null.String{},
		StripePlanID:  null.String{},
		AutoRenewal:   false,
		Status:        0,
		AppleSubsID:   null.String{},
		B2BLicenceID:  null.String{},
		AddOn:         addon.AddOn{},
		VIP:           false,
	}.Sync()
}
