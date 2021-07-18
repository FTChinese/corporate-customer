package order

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/guregu/null"
)

// CartItemBrief gives a brief statistics of CartItem.
// This is used when save all items of an order together in the's row.
type CartItemBrief struct {
	Price         price.Price `json:"price"`
	NewCopies     int64       `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int64       `json:"renewalCopies"` // How many renewals user purchased.
}

// CartItemSchema describes how a CartItem is saved in db.
// Each item in the shopping cart's items array is treated as a SQL row.
type CartItemSchema struct {
	ID      string
	OrderID string
	CartItem
	PriceOffPerCopy null.Float
}

type CartItem struct {
	Price     price.Price       `json:"price"`
	NewCopies int64             `json:"newCopies"`
	Renewals  []licence.Licence `json:"renewals"`
}

func (i CartItem) Brief() CartItemBrief {
	return CartItemBrief{
		Price:         i.Price,
		NewCopies:     i.NewCopies,
		RenewalCopies: int64(len(i.Renewals)),
	}
}

func (i CartItem) Schema(orderID string) CartItemSchema {
	return CartItemSchema{
		ID:              pkg.OrderItemID(),
		OrderID:         orderID,
		CartItem:        i,
		PriceOffPerCopy: null.Float{},
	}
}

type ShoppingCart struct {
	Items       []CartItem `json:"items"`
	ItemCount   int64      `json:"itemCount"`
	TotalAmount float64    `json:"totalAmount"`
}

func (s *ShoppingCart) ItemsBrief() []CartItemBrief {
	var items = make([]CartItemBrief, 0)
	for _, v := range s.Items {
		items = append(items, v.Brief())
	}

	return items
}
