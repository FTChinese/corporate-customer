// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
)

func (a Admin) ExistingExpLicence(p price.Price, to licence.Assignee) licence.ExpandedLicence {
	orderID := pkg.OrderID()
	return licence.ExpandedLicence{
		Licence:  licence.NewLicence(p, orderID, a.Creator()),
		Assignee: licence.AssigneeJSON{Assignee: to},
	}
}
