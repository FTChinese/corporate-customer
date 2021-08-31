package checkout

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/go-rest/chrono"
)

type OrderPaid struct {
	Payment
	Offers []PaymentOffer `json:"offers"`
}

func NewOrderPid(orderID string, params input.OrderPaidParams) OrderPaid {
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
// With LicenceQueue LEFT JOIN PaymentOffer
// using order_id, price_id and kind, you can get
// each licence's discount details.
type PaymentOffer struct {
	OrderID string `json:"orderId" db:"order_id"`
	Index   int64  `json:"-" db:"array_index"`
	input.PaymentOfferParams
}

func buildPaymentOffers(orderID string, offers []input.PaymentOfferParams) []PaymentOffer {
	var o = make([]PaymentOffer, 0)

	for i, v := range offers {
		o = append(o, PaymentOffer{
			OrderID:            orderID,
			Index:              int64(i),
			PaymentOfferParams: v,
		})
	}

	return o
}
