package input

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type PaymentParams struct {
	AmountPaid    float64           `json:"amountPaid" db:"amount_paid"`
	ApprovedBy    string            `json:"approvedBy" db:"approved_by"`
	Description   null.String       `json:"description" db:"description"`
	PaymentMethod pkg.PaymentMethod `json:"paymentMethod" db:"payment_method"`
	TransactionID null.String       `json:"transactionId" db:"transaction_id"`
}

type PaymentOfferParams struct {
	Copies          int64          `json:"copies" db:"copy_count"`
	Kind            enum.OrderKind `json:"kind" db:"kind"`
	Price           price.Price    `json:"price" db:"price"`
	PriceOffPerCopy float64        `json:"priceOffPerCopy" db:"price_off_per_copy"`
}

type OrderPaidParams struct {
	PaymentParams
	Offers []PaymentOfferParams `json:"offers"`
}
