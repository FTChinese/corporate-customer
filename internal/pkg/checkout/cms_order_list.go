package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
)

// CMSOrderRow is a row in the order table that is slightly
// different from the version used by admin.
type CMSOrderRow struct {
	Order
	Team TeamJSON `json:"team" db:"team"`
}

type CMSOrderList struct {
	pkg.PagedList
	Data []CMSOrderRow `json:"data"`
}
