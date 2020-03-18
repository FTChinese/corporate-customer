package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type Membership struct {
	SubsID     null.String    `json:"subsId" db:"subs_id"`
	Tier       enum.Tier      `json:"tier" db:"subs_tier"`
	Cycle      enum.Cycle     `json:"cycle" db:"subs_cycle"`
	ExpireDate chrono.Date    `json:"expireDate" db:"subs_expire_date"`
	PayMethod  enum.PayMethod `json:"payMethod" db:"subs_pay_method"`
	VIP        bool           `json:"vip" db:"is_vip"`
}
