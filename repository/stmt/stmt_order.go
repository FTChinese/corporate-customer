package stmt

import "github.com/FTChinese/b2b/models/sq"

const CreateOrder = `
INSERT INTO b2b.order (
	id,
	plan_id,
	discount_id,
	licence_id,
	team_id,
	checkout_id,
	amount,
	cycle_count,
	trial_days,
	kind,
	created_utc
) VALUES `

var OrderBuilder = sq.NewInsert().
	Into(sq.NewTable("b2b.order")).
	SetColumns([]sq.Column{
		sq.NewColumn("id"),
		sq.NewColumn("plan_id"),
		sq.NewColumn("discount_id"),
		sq.NewColumn("licence_id"),
		sq.NewColumn("team_id"),
		sq.NewColumn("checkout_id"),
		sq.NewColumn("amount"),
		sq.NewColumn("cycle_count"),
		sq.NewColumn("trial_days"),
		sq.NewColumn("kind"),
		sq.NewColumn("created_utc"),
	})

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
