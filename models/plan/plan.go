// Package plan describes the subscription plans we offered to
// b2b customers.
package plan

import (
	"github.com/FTChinese/go-rest/enum"
)

type BasePlan struct {
	PlanID    string     `json:"id" db:"plan_id"`
	Price     float64    `json:"price" db:"price"`
	Tier      enum.Tier  `json:"tier" db:"tier"`
	Cycle     enum.Cycle `json:"cycle" db:"cycle"`
	TrialDays int64      `json:"trialDays" db:"trial_days"`
}

// Plan is the output data structure.
// A plan may have discounts.
type Plan struct {
	BasePlan
	Discounts []Discount // null if discount does not exist when turned into JSON.
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

// DiscountPlansSchema is used as the scan target
// when retrieve a plans and its associated discounts
// in one shot. With LEFT JOIN, the rows retrieved
// is determined by the number of discounts.
type DiscountPlanSchema struct {
	BasePlan
	Discount
}

// ReduceDiscountPlan transforms a slice or DiscountPlanSchema
// retrieved from DB to Plan.
// For each DiscountPlanSchema the BasePlan part should
// be identical since they are the same one while the
// Discount might be different if there are more than 1
// discount record for this plan, or not exist if no discount
// is made for this plan.
func ReduceDiscountPlan(rows []DiscountPlanSchema) Plan {
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

// GroupedPlans is used to group discounts under each plan.
// The key is plan's id.
type GroupedPlans map[string]Plan

// GroupDiscountPlans is used to group Discounts
// into distinct plan.
// We used plan table LEFT JOIN discount table
// to retrieve a plan and its associated discounts.
// Therefore the left part of the result
// might have one plan row duplicated
// if a plan has multiple discounts.
// We need to group the plan into distinct ones
// and put the discounts under the Discounts
// field.
func GroupDiscountPlans(rows []DiscountPlanSchema) GroupedPlans {
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
