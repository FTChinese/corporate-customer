package checkout

const StmtInsertCartItem = `
INSERT INTO b2b.cart_item
SET id = :id,
	array_index = :array_index,
	order_id = :order_id,
	price = :price,
	new_copy_count = :new_copy_count,
	renewal_list = :renewal_list
`

// CartItemSchema is used to save/retrieve CartItem.
type CartItemSchema struct {
	ID         string `json:"id" db:"id"`
	OrderID    string `json:"order_d" db:"order_id"`
	ArrayIndex int64  `json:"-" db:"array_index"`
	CartItem
}
