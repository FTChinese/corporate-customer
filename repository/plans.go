package repository

import (
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/b2b/repository/stmt"
	"strings"
)

const stmtActivePlans = stmt.Plan + `
FROM subs.active_plan
ORDER BY tier ASC, cycle DESC`

func (env Env) LoadActivePlans() ([]plan.Plan, error) {
	var plans = make([]plan.Plan, 0)

	err := env.db.Select(&plans, stmtActivePlans)
	if err != nil {
		return plans, err
	}

	return plans, nil
}

const stmtActiveDiscounts = stmt.Discount + `
FROM subs.b2b_discount
WHERE FIND_IN_SET(plan_id, ?)
ORDER BY plan_id, quantity ASC`

// LoadDiscounts retrieves all discount schema for all currently
// active plans.
func (env Env) LoadActiveDiscounts(planIDs []string) ([]plan.Discount, error) {
	var discounts = make([]plan.Discount, 0)

	idList := strings.Join(planIDs, ",")

	err := env.db.Select(&discounts, stmtActiveDiscounts, idList)
	if err != nil {
		return discounts, err
	}

	return discounts, nil
}

const stmtPlan = stmt.Plan + `
FROM subs.plan
WHERE id = ?
LIMIT 1`

// LoadPlan retrieves a plan to which user is trying to subscribe.
func (env Env) LoadPlan(id string) (plan.Plan, error) {
	var p plan.Plan

	err := env.db.Get(&p, stmtPlan, id)

	if err != nil {
		return p, err
	}

	return p, nil
}

// order quantity by ascending order so that
// we doe not need to sort them in application.
const stmtDiscounts = stmt.Discount + `
FROM subs.b2b_discount
WHERE plan_id = ?
ORDER BY quantity ASC`

// LoadDiscounts retrieves all discounts under a plan.
func (env Env) LoadDiscounts(planID string) ([]plan.Discount, error) {
	var d = make([]plan.Discount, 0)

	err := env.db.Select(&d, stmtDiscounts, planID)
	if err != nil {
		return d, err
	}

	return d, nil
}
