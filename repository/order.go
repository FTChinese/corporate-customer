package repository

const stmtCreateOrder = `
INSERT into b2b.order
SET id = :order_id,
	plan_id = :plan_id,
	price = :price,
	amount = :amount,
	tier = :tier,
	cycle = :cycle,
	cycle_count = :cycle_count,
	trial_days = :trial_days,
	kind = :kind,
	created_utc = UTC_TIMESTAMP()`
