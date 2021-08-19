package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/guregu/null"
)

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

// CartItem describes an item user put into shopping cart.
// Usually user would have an array of CartItem submitted.
// This array is saved in two place: as an integral json value under
// an order row, and also save one by one with all details
// into a separate table.
type CartItem struct {
	Price     price.Price     `json:"price"`
	NewCopies int64           `json:"newCopies"`
	Renewals  RenewalLicences `json:"renewals"`
}

// CartItemSummary is an overview of CartItem.
type CartItemSummary struct {
	Price         price.Price `json:"price"`
	NewCopies     int64       `json:"newCopies"`     // How many new copies user purchased
	RenewalCopies int64       `json:"renewalCopies"` // How many renewals user purchased.
}

func newCartItemSummary(i CartItem) CartItemSummary {
	return CartItemSummary{
		Price:         i.Price,
		NewCopies:     i.NewCopies,
		RenewalCopies: int64(len(i.Renewals)),
	}
}

// CartItemSummaryList is a slice of CartItemSummary.
// It corresponds to the order.cart_items_summary field
// in DB so that when retrieving data, we don't need to
// retrieve all details of an order's items.
type CartItemSummaryList []CartItemSummary

func newCartItemSummaryList(items []CartItem) CartItemSummaryList {
	var list = make([]CartItemSummary, 0)
	for _, v := range items {
		list = append(list, newCartItemSummary(v))
	}

	return list
}

func (l CartItemSummaryList) Value() (driver.Value, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (l *CartItemSummaryList) Scan(src interface{}) error {
	if src == nil {
		*l = []CartItemSummary{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp []CartItemSummary
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*l = tmp
		return nil

	default:
		return errors.New("incompatible type to scan to []CartItemSummary")
	}
}

// OrderItem describes how a CartItem is saved in db.
// Each item in the shopping cart's items array is treated as a SQL row.
// Saved in b2b.order_item table with Price and Renewals fields as JSON.
type OrderItem struct {
	ID              string     `json:"id" db:"order_item_id"`
	OrderID         string     `json:"orderId" db:"order_id"`
	PriceOffPerCopy null.Float `json:"priceOffPerCopy" db:"price_off_per_copy"`
	CartItem
}

// NewOrderItem turns a CartItem into an OrderItem.
func NewOrderItem(orderID string, i CartItem) OrderItem {
	return OrderItem{
		ID:              pkg.OrderItemID(),
		OrderID:         orderID,
		CartItem:        i,
		PriceOffPerCopy: null.Float{},
	}
}
