package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/price"
)

// OrderItem is an overview of CartItem.
type OrderItem struct {
	Price         price.Price `json:"price"`
	NewCopies     int64       `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int64       `json:"renewalCopies"` // How many renewal user purchased.
}

// Turn an array of cart items into OrderItem.
func newOrderItemList(items []CartItem) OrderItemListJSON {
	var list = make([]OrderItem, 0)
	for _, v := range items {
		list = append(list, v.Summary())
	}

	return list
}