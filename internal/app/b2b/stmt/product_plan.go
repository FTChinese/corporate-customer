package stmt

// Products selects all products.
// Note the ORDER BY clause. The order is based
// on the order you define them.
const Products = `
SELECT id AS product_id,
	tier,
	heading,
	description,
	small_print,
	yearly_plan_id
FROM subs.product
ORDER BY tier ASC`

const planCols = `
p.id AS plan_id,
p.price AS price,
p.tier AS tier,
p.cycle AS cycle`

const discountCols = `
d.id AS discount_id,
IFNULL(d.quantity, 0) AS quantity,
IFNULL(d.price_off, 0) AS price_off`

// In the left join, discount table might be null,
// and to simplify things, the fields Discount type
// are not nullable types, so we use IFNULL to safe handle it.
const selectPlan = `
SELECT + ` + planCols + `,
` + discountCols + `
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
WHERE FIND_IN_SET(p.id, ?) > 0
ORDER BY p.tier ASC, p.cycle ASC, d.quantity ASC`
