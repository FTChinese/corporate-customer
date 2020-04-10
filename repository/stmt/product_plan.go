package stmt

const Products = `
SELECT id AS product_id,
	heading,
	description,
	small_print,
	tier,
	yearly_plan_id
FROM subs.product
ORDER BY tier ASC`

const planCols = `
p.id AS plan_id,
p.price AS price,
p.tier AS tier,
p.cycle AS cycle,
p.trial_days AS trial_days`

const selectPlan = `
SELECT p.id AS plan_id,
	p.price AS price,
	p.tier AS tier,
	p.cycle AS cycle,
	p.trial_days AS trial_days,
	d.quantity AS quantity,
	d.price_off AS price_off
FROM subs.plan AS p
	LEFT JOIN subs.b2b_discount AS d
	ON p.id = d.plan_id`

// Plan loads a single plan.
// The return is might contains multiple rows.
const Plan = selectPlan + `
WHERE p.id = ?
ORDER BY d.quantity ASC`

// PlansInSet loads multiple plans.
const ListPlans = selectPlan + `
WHERE FIND_IN_SET(p.id, ?)
ORDER BY p.tier ASC, p.cycle ASC, d.quantity ASC`
