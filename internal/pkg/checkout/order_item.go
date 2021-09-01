package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/price"
)

// OrderItem is a summary of CartItem.
// It will be saved as an array into one of order columns.
type OrderItem struct {
	Price         price.Price `json:"price"`
	NewCopies     int         `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int         `json:"renewalCopies"` // How many renewal user purchased.
}

// Turn an array of cart items into OrderItem.
func newOrderItemList(items []CartItem) OrderItemListJSON {
	var list = make([]OrderItem, 0)
	for _, v := range items {
		list = append(list, v.OrderItem())
	}

	return list
}
