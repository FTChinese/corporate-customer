package plan

import "github.com/guregu/null"

// Discount is the amount to subtract from Plan.Price when
// user purchases licence in bulk.
// For example, when purchase 10 copies, the price of each
// licence will off 10.
// Example:
type Discount struct {
	DiscountID null.Int `json:"id" db:"discount_id"`
	Quantity   int64    `json:"quantity" db:"quantity"`  // The amount of minimum copies of licences purchased when this discount becomes available.
	PriceOff   float64  `json:"priceOff" db:"price_off"` // Deducted from Plan.Price
}

// IsZero tests if a the discount actually exist since
// we use a JOIN to retrieve the plan together with its
// discount in one go.
func (d Discount) IsZero() bool {
	return d.PriceOff <= 0
}

// DiscountPlansSchema contains a discount schema and its plan.
// This is used as the scan target when retrieved plan and its discount
// from DB one one shot, using LEFT JOIN. The rows retrieved
// is determined by the number of discounts and the plan row might be
// duplicated.
type DiscountPlan struct {
	BasePlan
	Discount
}

// PayableAmount calculates how much use should pay for this
// play after applicable discount amount is deducted from
// this plan's price.
func (d DiscountPlan) PayableAmount() float64 {
	return d.Price - d.PriceOff
}
