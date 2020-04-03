package stmt

const Plan = `
SELECT id AS plan_id,
	price,
	tier,
	cycle,
	trial_days,
	created_utc`

const Discount = `
SELECT plan_id,
	quantity,
	price_off`

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

const Licence = `
SELECT id AS licence_id,
	plan_id,
	team_id,
	expire_date,
	assignee_id
	is_active,
	created_utc,
	updated_utc
FROM b2b.licence`

// Invitation is used to retrieve to list all, show a single one to admin, an expanded
// version when user is accepting it, or
// locking when granting licence.
const Invitation = `
SELECT id AS invitation_id,
	licence_id,
	team_id,
	LOWER(HEX(token)) AS token,
	invitee_email,
	description,
	accepted,
	revoked,
	created_utc,
	updated_utc
FROM b2b.invitation`
