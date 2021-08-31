package checkout

import "github.com/FTChinese/ftacademy/internal/pkg/reader"

type MembershipRenewed struct {
	Latest  reader.Membership
	Archive reader.MemberSnapshot
}
