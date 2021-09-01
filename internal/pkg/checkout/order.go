package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
)

// OrderItem is a summary of CartItem.
// It will be saved as an array into one of order columns.
type OrderItem struct {
	Price         price.Price `json:"price"`
	NewCopies     int         `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int         `json:"renewalCopies"` // How many renewal user purchased.
}

// Order is what a shopping cart should create.
type Order struct {
	ID string `json:"id" db:"order_id"`
	admin.Creator
	AmountPayable float64           `json:"amountPayable" db:"amount_payable"`
	CreatedUTC    chrono.Time       `json:"createdUtc" db:"created_utc"`
	ItemCount     int64             `json:"itemCount" db:"item_count"`
	ItemList      OrderItemListJSON `json:"itemList" db:"item_list"`
	Status        Status            `json:"status" db:"current_status"`
}

func NewOrder(cart ShoppingCart, p admin.PassportClaims) Order {
	return Order{
		ID: pkg.OrderID(),
		Creator: admin.Creator{
			AdminID: p.AdminID,
			TeamID:  p.TeamID.String,
		},
		AmountPayable: cart.TotalAmount,
		CreatedUTC:    chrono.TimeUTCNow(),
		ItemList:      cart.OrderItemList(),
		ItemCount:     cart.ItemCount,
		Status:        StatusPending,
	}
}

func (o Order) ChangeStatus(s Status) Order {
	o.Status = s

	return o
}

func (o Order) IsFinal() bool {
	return o.Status == StatusPaid || o.Status == StatusCancelled
}

// OrderList contains a list of orders
type OrderList struct {
	pkg.PagedList
	Data []Order `json:"data"`
}

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
