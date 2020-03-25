package stmt

const CreateOrder = `
INSERT into b2b.transaction
SET id = :tnx_id,
	price = :price,
	amount = :amount,
	tier = :tier,
	cycle = :cycle,
	cycle_count = :cycle_count,
	trial_days = :trial_days,
	kind = :kind,
	created_utc = UTC_TIMESTAMP()`
