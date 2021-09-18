package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/go-rest/chrono"
)

// ShoppingCart is used to hold data submitted by client.
// When storing a shopping cart, its data is decomposed:
// A summary of items is saved in an order row;
// the items array is saved into a separate table one by one.
// Each licence derived from Items will be saved into
// licence_queue table
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

type OrderSchemaBuilder struct {
	orderID string
	cart    ShoppingCart
	creator admin.Creator
}

func NewOrderSchemaBuilder(cart ShoppingCart, pp admin.PassportClaims) OrderSchemaBuilder {
	return OrderSchemaBuilder{
		orderID: pkg.OrderID(),
		cart:    cart,
		creator: admin.Creator{
			AdminID: pp.AdminID,
			TeamID:  pp.TeamID.String,
		},
	}
}

func (b OrderSchemaBuilder) Order() Order {
	return Order{
		ID:            b.orderID,
		Creator:       b.creator,
		AmountPayable: b.cart.TotalAmount,
		CreatedUTC:    chrono.TimeUTCNow(),
		ItemList:      b.cart.OrderItemList(),
		ItemCount:     b.cart.ItemCount,
		Status:        StatusPending,
	}
}

func (b OrderSchemaBuilder) CartItemSchema() []CartItemSchema {
	var s = make([]CartItemSchema, 0)

	for _, item := range b.cart.Items {
		s = append(s, item.Schema(b.orderID, b.creator))
	}

	return s
}

func (b OrderSchemaBuilder) TransactionList() []LicenceTransaction {
	var txnList = make([]LicenceTransaction, 0)

	// Loop over each cart item
	for _, item := range b.cart.Items {
		// Loop over each new copy of a cart item.
		for i := 0; i < int(item.NewCopies); i++ {
			txn := NewLicenceTransaction(
				b.orderID,
				item.Price,
				b.creator,
				licence.ExpandedLicence{},
			)
			txnList = append(txnList, txn)
		}

		// Loop over each renewed licence.
		for _, lic := range item.Renewals {
			item := NewLicenceTransaction(
				b.orderID,
				item.Price,
				b.creator,
				lic,
			)
			txnList = append(txnList, item)
		}
	}

	return txnList
}

func (b OrderSchemaBuilder) Build() OrderInputSchema {
	return OrderInputSchema{
		OrderRow:     b.Order(),
		CartItems:    b.CartItemSchema(),
		Transactions: b.TransactionList(),
	}
}
