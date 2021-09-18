package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
)

type OrderPaid struct {
	Payment
	Offers []PaymentOffer `json:"offers"`
}

func NewOrderPaid(orderID string, params input.OrderPaidParams) OrderPaid {
	return OrderPaid{
		Payment: Payment{
			OrderID:       orderID,
			ApprovedUTC:   chrono.TimeNow(),
			PaymentParams: params.PaymentParams,
		},
		Offers: buildPaymentOffers(orderID, params.Offers),
	}
}

// Payment describes the details of an order's payment.
type Payment struct {
	OrderID     string      `json:"orderId" db:"order_id"`
	ApprovedUTC chrono.Time `json:"approvedUtc" db:"approved_utc"`
	input.PaymentParams
}

// PaymentOffer describes how discount is used for a
// price of specific kind.
// With LicenceTransaction LEFT JOIN PaymentOffer
// using order_id, price_id and kind, you can get
// each licence's discount details.
type PaymentOffer struct {
	OrderID string `json:"orderId" db:"order_id"`
	input.PaymentOfferParams
}

func buildPaymentOffers(orderID string, offers []input.PaymentOfferParams) []PaymentOffer {
	var o = make([]PaymentOffer, 0)

	for _, v := range offers {
		o = append(o, PaymentOffer{
			OrderID:            orderID,
			PaymentOfferParams: v,
		})
	}

	return o
}

// PriceOfLicenceTxn maps the LicenceTransaction id to the price it used.
type PriceOfLicenceTxn struct {
	TxnID string      `db:"txn_id"`
	Price price.Price `db:"price"`
}

type PaymentError struct {
	TxnID      string      `db:"txn_id"`
	Message    string      `db:"error_message"`
	CreatedUTC chrono.Time `db:"created_utc"`
}

func NewPaymentError(txnID string, err error) PaymentError {
	return PaymentError{
		TxnID:      txnID,
		Message:    err.Error(),
		CreatedUTC: chrono.TimeNow(),
	}
}
