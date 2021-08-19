package checkout

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
