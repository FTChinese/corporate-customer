package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/ftacademy/pkg/sq"
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

func (s *ShoppingCart) ProductsBrief() []CheckoutProduct {
	var items = make([]CheckoutProduct, 0)
	for _, v := range s.Items {
		items = append(items, v.ProductBrief())
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

// ProductBrief gives an overview of a product put in cart.
func (i CartItem) ProductBrief() CheckoutProduct {
	return CheckoutProduct{
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

// CartItemSchema describes how a CartItem is saved in db.
// Each item in the shopping cart's items array is treated as a SQL row.
type CartItemSchema struct {
	ID              string     `json:"id" db:"order_item_id"`
	OrderID         string     `json:"orderId" db:"order_id"`
	PriceOffPerCopy null.Float `json:"priceOffPerCopy" db:"price_off_per_copy"`
	CartItem
}

func (s CartItemSchema) RowValues() []interface{} {
	return []interface{}{
		s.ID,
		s.OrderID,
		s.CartItem,
		s.PriceOffPerCopy,
	}
}

// RenewalLicences to ease retrieve/save an array of
// licences into db.
type RenewalLicences []licence.Licence

func (rl RenewalLicences) Value() (driver.Value, error) {
	j, err := json.Marshal(rl)
	if err != nil {
		return nil, err
	}

	return string(j), nil
}

func (rl *RenewalLicences) Scan(src interface{}) error {
	if src == nil {
		*rl = []licence.Licence{}
		return nil
	}
	switch s := src.(type) {
	case []byte:
		var tmp []licence.Licence
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*rl = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to []Licence")
	}
}

type CartItems []CartItemSchema

// Each implements Enumerable interface.
// Usage: `BuildInsertValues(CartItems)...`
func (ci CartItems) Each(handler func(row sq.InsertRow)) {
	for _, c := range ci {
		handler(c)
	}
}
