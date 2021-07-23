// +build !production

package price

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

var MockPriceStdYear = Price{
	ID: "plan_MynUQDQY1TSQ",
	Edition: Edition{
		Tier:  enum.TierStandard,
		Cycle: enum.CycleYear,
	},
	Active:     true,
	Currency:   CurrencyCNY,
	LiveMode:   true,
	Nickname:   null.String{},
	ProductID:  "prod_zjWdiTUpDN8l",
	Source:     SourceFTC,
	UnitAmount: 298,
}

var MockPricePrm = Price{
	ID: "plan_vRUzRQ3aglea",
	Edition: Edition{
		Tier:  enum.TierPremium,
		Cycle: enum.CycleYear,
	},
	Active:     true,
	Currency:   CurrencyCNY,
	LiveMode:   true,
	Nickname:   null.String{},
	ProductID:  "prod_IaoK5SbK79g8",
	Source:     SourceFTC,
	UnitAmount: 1998,
}
