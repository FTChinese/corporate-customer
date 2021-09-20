package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
)

// CartItem describes an item user put into shopping cart.
// Usually user would have an array of CartItem submitted.
// This array is saved in two place: as an integral json value under
// an order row, and also save one by one with all details
// into a separate table.
type CartItem struct {
	Price     price.Price        `json:"price" db:"price"`
	NewCopies int64              `json:"newCopies" db:"new_copy_count"`
	Renewals  ExpLicenceListJSON `json:"renewals" db:"renewal_list"` // This field is saved as SQL JSON type.
}

func (ci CartItem) OrderItem() OrderItem {
	return OrderItem{
		Price:         ci.Price,
		NewCopies:     int(ci.NewCopies),
		RenewalCopies: len(ci.Renewals),
	}
}

func (ci CartItem) Schema(orderID string, creator admin.Creator) CartItemSchema {
	return CartItemSchema{
		OrderID:    orderID,
		CartItem:   ci,
		Creator:    creator,
		CreatedUTC: chrono.TimeNow(),
	}
}

// CartItemSchema is used to save/retrieve CartItem into
// table order_item
type CartItemSchema struct {
	OrderID string `json:"orderId" db:"order_id"`
	CartItem
	admin.Creator
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
}
