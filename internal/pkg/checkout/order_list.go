package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/go-rest/chrono"
)

// OrderRow describes the details of a transaction for a
// session of shopping.
// If a transaction is used to purchase a new licence, the
// licence should be created together with the order but marked
// as inactive. Once the transaction is confirmed,
// the licence will be activated and the admin is allowed to
// invite someone to use this licence.
// If a transaction is used to renew/upgrade a licence,
// the licence associated with it won't be touched until
// it is confirmed, which will result licence extended or
// upgraded and the membership (if the licence is granted
// to someone) will be backed up and updated corresponding.
type OrderRow struct {
	BaseOrder
	// An array of products, together with the quantities, use is trying to purchase.
	ItemSummaryList CartItemSummaryList `json:"items" db:"cart_items_summary"`
}

func NewOrderRow(cart ShoppingCart, p admin.PassportClaims) OrderRow {
	return OrderRow{
		BaseOrder: BaseOrder{
			ID:            pkg.OrderID(),
			AmountPayable: cart.TotalAmount,
			CreatedBy:     p.AdminID,
			CreatedUTC:    chrono.TimeUTCNow(),
			ItemCount:     cart.ItemCount,
			Status:        StatusPending,
			TeamID:        p.TeamID.String,
		},
		ItemSummaryList: newCartItemSummaryList(cart.Items),
	}
}

// OrderRowList contains a list of orders
type OrderRowList struct {
	pkg.PagedList
	Data []OrderRow `json:"data"`
}
