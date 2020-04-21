// Package plan describes the subscription plans we offered to
// b2b customers.
package plan

import (
	"github.com/FTChinese/go-rest/enum"
)

type BasePlan struct {
	PlanID string     `json:"id" db:"plan_id"`
	Price  float64    `json:"price" db:"price"`
	Tier   enum.Tier  `json:"tier" db:"tier"`
	Cycle  enum.Cycle `json:"cycle" db:"cycle"`
}

// Plan is the output data structure.
// A plan may have discounts.
type Plan struct {
	BasePlan
	Discounts []Discount `json:"discounts"` // null if discount does not exist when turned into JSON.
}

// NewPlan transforms a slice or DiscountPlan
// retrieved from DB to Plan.
// For example, you have a plan with two discount rows.
// The retrieved data looks like:
//
//      plan_id       |  price  |    tier    |  cycle | trial | quantity | price_off
// --------------------------------------------------------------------------------
// "plan_ICMPPM0UXcpZ"	"258.00"	"standard"	"year"	"3"	   "10"	    "15.00"
// "plan_ICMPPM0UXcpZ"	"258.00"	"standard"	"year"	"3"    "20"	    "25.00"
//
// The BasePlan part are identical. Use any one of them. The discounts part are
// different.
// You also need to take into account the fact that a plan might have no discount
// at all. In such case, you will only get one row with quantity and price_off set
// to zero. You should ignore the discount part, which is handled by the
// AddDiscount method.
func NewPlan(rows []DiscountPlan) Plan {
	if len(rows) == 0 {
		return Plan{}
	}

	// There must be multiple discounts under this plan.
	p := Plan{
		BasePlan: rows[0].BasePlan, // Use any rows's BasePlan works since they are identical.
	}

	for _, v := range rows {
		// Just a a precaution.
		if v.PlanID != p.PlanID {
			continue
		}
		p.AddDiscount(v.Discount)
	}

	return p
}

// AddDiscount appends a discount to this plan.
// Zero value is discarded.
func (p *Plan) AddDiscount(d Discount) {
	if d.IsZero() {
		return
	}
	if p.Discounts == nil {
		p.Discounts = make([]Discount, 0)
	}

	p.Discounts = append(p.Discounts, d)
}

// FindDiscount find out which discount will be used
// for a copies of q licences.
// This problem could be simplified to find out which range
// of a number falls into among several ascending-ordered
// numbers.
// For example, if this plan have those discounts:

// Quantity | Price off
// --------------------
// 10       |  19
// 20       |  29
// 30       |  39
// ----------------------
// If user purchased 5 copies, no discount;
// 15 copies fall into the first tier;
// 25 copies fall into the second tier;
// above 30 copies, always use the third tier.
func (p Plan) FindDiscount(q int64) Discount {
	if p.Discounts == nil || len(p.Discounts) == 0 {
		return Discount{}
	}

	// Use a zero value to init.
	var previous = Discount{}
	for _, v := range p.Discounts {
		if q > previous.Quantity && q < v.Quantity {
			return previous
		}

		previous = v
	}
	// After the loop, previous is the last element in the slice.
	// This indicates the quantity is larger than the last
	// discountable tier.
	return previous
}

// DiscountPlan produces a DiscountPlan used to record
// the plan and discount details upon checkout.
// The data is saved as a JSON document for reference only.
func (p Plan) DiscountPlan(q int64) DiscountPlan {
	d := p.FindDiscount(q)
	return DiscountPlan{
		BasePlan: p.BasePlan,
		Discount: d,
	}
}

// GroupedPlans is used to group discounts under each plan.
// The key is plan's id.
type GroupedPlans map[string]Plan

// NewGroupedPlans is used to group Discounts
// into distinct plan.
// We used plan table LEFT JOIN discount table
// to retrieve a plan and its associated discounts.
// Therefore the left part of the result
// might have one plan row duplicated
// if a plan has multiple discounts.
// We need to group the plan into distinct ones
// and put the discounts under the Discounts
// field.
func NewGroupedPlans(rows []DiscountPlan) GroupedPlans {
	var plans = make(GroupedPlans)

	for _, v := range rows {
		plan, ok := plans[v.PlanID]
		// If v is no put into the map.
		if !ok {
			plan = Plan{
				BasePlan:  v.BasePlan,
				Discounts: make([]Discount, 0),
			}
		}

		plan.AddDiscount(v.Discount)

		plans[v.PlanID] = plan
	}

	return plans
}
