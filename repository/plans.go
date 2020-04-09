package repository

import (
	"github.com/FTChinese/b2b/models/plan"
	"github.com/FTChinese/b2b/repository/stmt"
	"strings"
)

// PlansInSet loads a list of active plans
// and returns a map with plan id as key.
func (env Env) PlansInSet(planIDs []string) (plan.Plans, error) {
	idSet := strings.Join(planIDs, ",")
	var raws = make([]plan.DiscountPlanSchema, 0)

	err := env.db.Select(&raws, stmt.ListPlans, idSet)
	if err != nil {
		return plan.Plans{}, err
	}

	return plan.GroupRawPlans(raws), nil
}

// LoadPlan retrieves a single plan.
func (env Env) LoadPlan(id string) (plan.Plan, error) {
	var raws []plan.DiscountPlanSchema

	err := env.db.Get(&raws, stmt.Plan, id)

	if err != nil {
		return plan.Plan{}, err
	}

	return plan.BuildPlan(raws), nil
}
