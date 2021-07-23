package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
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
// Used for JSON output only.
type Order struct {
	BaseOrder
	Items []OrderItem `json:"items"`
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

// NewOrderItem turns a CartItem into an OrderItem.
func NewOrderItem(orderID string, i CartItem) OrderItem {
	return OrderItem{
		ID:              pkg.OrderItemID(),
		OrderID:         orderID,
		CartItem:        i,
		PriceOffPerCopy: null.Float{},
	}
}
