// Package plan describes the subscription plans we offered to
// b2b customers.
package plan

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

type Plan struct {
	ID         string      `db:"plan_id"`
	Price      float64     `db:"price"`
	Tier       enum.Tier   `db:"tier"`
	Cycle      enum.Cycle  `db:"cycle"`
	TrialDays  int64       `db:"trial_days"`
	CreatedUTC chrono.Time `db:"created_utc"`
	Order      int64       `db:"order_weight"`
	Discounts  []Discount
}

// FindDiscount find out which discount will be used
// for a copies of q licences.
// This problem could be simplified to find out which range
// of a number falls into among several ascending-ordered
// numbers.
// For example, if there are three tiers:
// 10 copies for price off 19
// 20 copies for price off 29
// 30 copies for price off 39
// above 30 copies use the 30 tier.
// Given a purchase of 25 copies, we should use the 20 tier;
// for 40 copies, use the 30 tier.
// The Discounts array should be sorted by Quantity on
// ascending order.
// The final price payable: p.Price - Discount.PriceOff
func (p Plan) FindDiscount(q int64) Discount {
	if p.Discounts == nil {
		return Discount{}
	}

	if q < p.Discounts[0].Quantity {
		return p.Discounts[0]
	}

	l := len(p.Discounts)
	if q > p.Discounts[l-1].Quantity {
		return p.Discounts[l-1]
	}

	// Use a zero value to init.
	var previous = Discount{}
	for _, v := range p.Discounts {
		if q > previous.Quantity && q < v.Quantity {
			return previous
		}

		previous = v
	}

	return Discount{}
}

// Discount is the amount to subtract from Plan.Price when
// user purchases licence in bulk.
// For example, when purchase 10 copies, the price of each
// licence will off 10.
// Example:
type Discount struct {
	PlanID   string  `db:"plan_id"`
	Quantity int64   `db:"quantity"`  // The amount of minimum copies of licences purchased when this discount becomes available.
	PriceOff float64 `db:"price_off"` // Deducted from Plan.Price
}

func (d Discount) IsZero() bool {
	return d.PriceOff == 0
}
