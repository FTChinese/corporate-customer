package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
)

// ShoppingCart is used to hold data submitted by client.
// When storing a shopping cart, its data is decomposed:
// A summary of items is saved in an order row;
// the items array is saved into a separate table one by one.
// Table `order` to `order_item` has a one-many relation.
type ShoppingCart struct {
	Items       []CartItem `json:"items"`
	ItemCount   int64      `json:"itemCount"`
	TotalAmount float64    `json:"totalAmount"`
}

// OrderItemList turns Items field to a list to be saved
// as a column of order row.
func (c ShoppingCart) OrderItemList() OrderItemListJSON {
	var list = make([]OrderItem, 0)
	for _, v := range c.Items {
		list = append(list, v.OrderItem())
	}

	return list
}

// CartItemSchema turns Items to a list to be saved one
// row per element in a table separate from order.
func (c ShoppingCart) CartItemSchema(orderID string) []CartItemSchema {
	var s = make([]CartItemSchema, 0)

	for _, item := range c.Items {
		s = append(s, item.Schema(orderID))
	}

	return s
}

// LicenceQueue create db schema to queue the licence-to-be-created.
func (c ShoppingCart) LicenceQueue(orderID string) []LicenceQueue {
	var queue = make([]LicenceQueue, 0)

	// Loop over each cart item
	index := 0
	for _, item := range c.Items {
		// Loop over each new copy of a cart item.
		for i := 0; i < int(item.NewCopies); i++ {
			index += 1

			qi := NewLicenceQueue(orderID, item.Price, licence.ExpandedLicence{}, index)
			queue = append(queue, qi)
		}

		// Loop over each renewed licence.
		for _, lic := range item.Renewals {
			index += 1

			item := NewLicenceQueue(orderID, item.Price, lic, index)
			queue = append(queue, item)
		}
	}

	return queue
}
