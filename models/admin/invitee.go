package admin

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Invitee is a member of a team who will be granted
// a licence.
// When an invitee received an email, it should click the link
// in the email.
// When the link is opened, we take the following steps in order:
// 1. First check the if the token exists. A token is valid only
// when it does exist, is not accepted yet, and not revoked,
//and not expired.
// 2. Then we should check if the licence to be granted is
// still available.
// 3. Check if the invitee has an account at FTC. If not, ask it to sign up.
// 4. For existing user, check whether the account has a valid
// membership with it. It it does have one, deny the granting.
// 5. Link the licence to user's ftc id;
// 6. Insert membership if user does not have membership yet, or
// backup existing membership and update membership.
// 6. Mark the invitation as accepted;
type Invitee struct {
	FtcID      string         `db:"ftc_id"`
	Email      string         `db:"email"`
	MemberID   null.String    `db:"mmb_id"`
	Tier       enum.Tier      `db:"mmb_tier"`
	Cycle      enum.Cycle     `db:"mmb_cycle"`
	ExpireDate chrono.Date    `db:"mmb_expire_date"`
	PayMethod  enum.PayMethod `db:"mmb_pay_method"`
	VIP        bool           `db:"is_vip"`
}

func (i Invitee) HasMembership() bool {
	return i.Tier != enum.TierNull
}
