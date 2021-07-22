package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/pkg/sq"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type BaseOrder struct {
	ID            string      `json:"id" db:"order_id"`
	AmountPayable float64     `json:"amountPayable" db:"amount_payable"`
	CreatedBy     string      `json:"createdBy" db:"created_by"`
	CreatedUTC    chrono.Time `json:"createdUtc" db:"created_utc"`
	ItemCount     int64       `json:"itemCount" db:"item_count"`
	Status        Status      `json:"status" db:"current_status"`
	TeamID        string      `json:"teamId" db:"team_id"`
}

// Payment describes the details the an order's payment.
type Payment struct {
	AmountPaid    null.Float    `json:"amountPaid" db:"amount_paid"`
	ApprovedBy    null.String   `json:"approvedBy" db:"approved_by"`
	ApprovedUTC   chrono.Time   `json:"approvedUtc" db:"approved_utc"`
	Description   null.String   `json:"description" db:"description"`
	PaymentMethod PaymentMethod `json:"paymentMethod" db:"payment_method"`
	TransactionID null.String   `json:"transactionId" db:"transaction_id"`
}

// Order contains all details of what user wanted to buy,
// how payment is handled.
type Order struct {
	BaseOrder
	CartItems []OrderItem `json:"cartItems"`
	Payment
}

// OrderItem describes how a CartItem is saved in db.
// Each item in the shopping cart's items array is treated as a SQL row.
type OrderItem struct {
	ID              string     `json:"id" db:"order_item_id"`
	OrderID         string     `json:"orderId" db:"order_id"`
	PriceOffPerCopy null.Float `json:"priceOffPerCopy" db:"price_off_per_copy"`
	CartItem
}

func (s OrderItem) RowValues() []interface{} {
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

type CartItems []OrderItem

// Each implements Enumerable interface.
// Usage: `BuildInsertValues(CartItems)...`
func (ci CartItems) Each(handler func(row sq.InsertRow)) {
	for _, c := range ci {
		handler(c)
	}
}
