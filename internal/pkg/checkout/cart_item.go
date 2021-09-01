package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/price"
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

func (ci CartItem) Schema(orderID string, i int) CartItemSchema {
	return CartItemSchema{
		ID:         pkg.CartItemID(),
		OrderID:    orderID,
		ArrayIndex: int64(i),
		CartItem:   ci,
	}
}

// CartItemSchema is used to save/retrieve CartItem.
type CartItemSchema struct {
	ID         string `json:"id" db:"id"`
	OrderID    string `json:"order_d" db:"order_id"`
	ArrayIndex int64  `json:"-" db:"array_index"`
	CartItem
}
