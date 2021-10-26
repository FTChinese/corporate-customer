package reader

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type BaseAccount struct {
	FtcID        string      `json:"id"`       // FTC's uuid
	UnionID      null.String `json:"unionId"`  // Wechat's union id
	StripeID     null.String `json:"stripeId"` // Stripe's id
	Email        string      `json:"email"`    // Required, unique. Max 64.
	Mobile       null.String `json:"mobile"`
	UserName     null.String `json:"userName"` // Required, unique. Max 64.
	AvatarURL    null.String `json:"avatarUrl"`
	IsVerified   bool        `json:"isVerified"`
	CampaignCode null.String `json:"campaignCode"`
}

// Wechat contain the essential data to identify a wechat user.
type Wechat struct {
	WxNickname  null.String `json:"nickname"`
	WxAvatarURL null.String `json:"avatarUrl"`
}

type Account struct {
	BaseAccount
	LoginMethod enum.LoginMethod `json:"loginMethod"`
	Wechat      Wechat           `json:"wechat"`
	Membership  Membership       `json:"membership"`
}
