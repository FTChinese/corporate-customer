package plan

// GetCartPlansIDs collects the plan ids of items
// in a cart and use it to retrieve the plans
// using FIND_IN_SET
func GetCartPlanIDs(items []CartItem) []string {
	var ids = make([]string, 0)

	for _, item := range items {
		ids = append(ids, item.PlanID)
	}

	return ids
}

// CartItem is the plan user is subscribing
// and the number of copies for this plan.
type CartItem struct {
	PlanID   string `json:"planId"`
	Quantity int64  `json:"quantity"`
}
