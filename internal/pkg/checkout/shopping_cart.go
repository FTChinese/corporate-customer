package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/guregu/null"
)

// ShoppingCart is used to hold data submitted by client.
// When storing a shopping cart, its data is decomposed:
// A summary of items is save in an order row;
// the items array is save into a separate table one by one.
// Table `order` to `order_item` has a one-many relation.
type ShoppingCart struct {
	Items       []CartItem `json:"items"`
	ItemCount   int64      `json:"itemCount"`
	TotalAmount float64    `json:"totalAmount"`
}

func (s ShoppingCart) ProductsBrief() []OrderedProduct {
	var items = make([]OrderedProduct, 0)
	for _, v := range s.Items {
		items = append(items, v.ProductOrdered())
	}

	return items
}

func (s ShoppingCart) ItemRows(orderID string) []OrderItem {
	var items = make([]OrderItem, 0)
	for _, v := range s.Items {
		items = append(items, v.Schema(orderID))
	}

	return items
}

// CartItem describes how many copies user purchased a product.
// Saved in b2b.order_item table with Price and Renewals fields as JSON.
type CartItem struct {
	Price     price.Price     `json:"price" db:"price_snapshot"`
	NewCopies int64           `json:"newCopies" db:"new_copy_count"`
	Renewals  RenewalLicences `json:"renewals" db:"renewal_licences"`
}

// ProductOrdered gives an overview of a product put in cart.
func (i CartItem) ProductOrdered() OrderedProduct {
	return OrderedProduct{
		Price:         i.Price,
		NewCopies:     i.NewCopies,
		RenewalCopies: int64(len(i.Renewals)),
	}
}

func (i CartItem) Schema(orderID string) OrderItem {
	return OrderItem{
		ID:              pkg.OrderItemID(),
		OrderID:         orderID,
		CartItem:        i,
		PriceOffPerCopy: null.Float{},
	}
}
