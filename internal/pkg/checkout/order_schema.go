package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
)

// OrderInputSchema is db schema to save a shopping cart.
// A shopping cart is dissected into 3 tables.
type OrderInputSchema struct {
	OrderRow  Order
	CartItems []CartItemSchema
	Queue     []LicenceQueue // All copies created from shopping cart items.
}

func NewOrderInputSchema(cart ShoppingCart, p admin.PassportClaims) OrderInputSchema {
	o := NewOrder(cart, p)

	return OrderInputSchema{
		OrderRow:  o,
		CartItems: cart.CartItemSchema(o.ID),
		Queue:     cart.LicenceQueue(o.ID),
	}
}
