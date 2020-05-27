package products

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	"github.com/FTChinese/ftacademy/internal/pkg/plan"
	"strings"
)

// PlansInSet loads a list of active plans
// and returns a map with plan id as key.
func (env Env) PlansInSet(planIDs []string) (plan.GroupedPlans, error) {
	// TODO: find in cache. If any of the them is not found in cache, then retrieve all fro DB
	idSet := strings.Join(planIDs, ",")
	var raws = make([]plan.DiscountPlan, 0)

	err := env.db.Select(&raws, stmt.ListPlans, idSet)
	if err != nil {
		return plan.GroupedPlans{}, err
	}

	// Cache each plan
	return plan.NewGroupedPlans(raws), nil
}

// LoadPlan retrieves a single plan.
func (env Env) LoadPlan(id string) (plan.Plan, error) {
	// TODO: load from cache first.
	var rows []plan.DiscountPlan

	err := env.db.Get(&rows, stmt.Plan, id)

	if err != nil {
		return plan.Plan{}, err
	}

	// TODO: cache this plan.
	return plan.NewPlan(rows), nil
}
