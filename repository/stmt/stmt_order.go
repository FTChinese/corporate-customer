package stmt

const Order = `
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
