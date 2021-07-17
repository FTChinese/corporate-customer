package order

// Deprecated
const selectOrder = `
SELECT id AS order_id,
	plan_id,
	licence_id,
	team_id,
	amount,
	cycle_count,
	period_start,
	period_end,
	kind,
	created_utc,
	confirmed_utc
FROM b2b.order`

// Deprecated
const GetOrder = selectOrder + `
WHERE id = ?
	AND team_id = ?
LIMIT 1`

const ListOrder = selectOrder + `
WHERE id = ?
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const CountOrder = `
SELECT COUNT(*)
FROM b2b.order
WHERE team_id = ?`
