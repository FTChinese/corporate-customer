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

// CartItemSchema creates db schema to save items array.
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
	for _, item := range c.Items {
		// Loop over each new copy of a cart item.
		for i := 0; i < int(item.NewCopies); i++ {
			qi := NewLicenceQueue(orderID, item.Price, licence.ExpandedLicence{}, i)
			queue = append(queue, qi)
		}

		// Loop over each renewed licence.
		for i, lic := range item.Renewals {
			item := NewLicenceQueue(orderID, item.Price, lic, i)
			queue = append(queue, item)
		}
	}

	return queue
}
