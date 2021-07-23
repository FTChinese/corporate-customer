package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/sq"
)

// OrderInputSchema is used to hold data to be saved for
// a shopping cart.
// A shopping cart is dissected into two parts:
// * a single base row for each shopping cart;
// * multiple rows for each product put into the cart.
type OrderInputSchema struct {
	OrderRow OrderRow
	ItemRows []OrderItem
}

func NewOrderInputSchema(cart ShoppingCart, p admin.PassportClaims) OrderInputSchema {
	bo := NewOrderRow(cart, p)

	var oi = make([]OrderItem, 0)
	for _, v := range cart.Items {
		oi = append(oi, NewOrderItem(bo.ID, v))
	}

	return OrderInputSchema{
		OrderRow: bo,
		ItemRows: oi,
	}
}

// RowValues build the a row of values for SQL bulk insert.
// Used together with Enumerable.Each to build the VALUES
// part of SQL bulk insert.
// These are experimental features that might never be used.
func (s OrderItem) RowValues() []interface{} {
	return []interface{}{
		s.ID,
		s.OrderID,
		s.CartItem,
		s.PriceOffPerCopy,
	}
}

// OrderItems is used to implement Enumerable interface.
type OrderItems []OrderItem

// Each implements Enumerable interface to loop over each
// row of bulk insert.
// Usage: `BuildInsertValues(OrderItems)...`
func (ci OrderItems) Each(handler func(row sq.InsertRow)) {
	for _, c := range ci {
		handler(c)
	}
}
