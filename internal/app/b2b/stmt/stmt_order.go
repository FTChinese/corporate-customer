package stmt

import (
	sq2 "github.com/FTChinese/ftacademy/pkg/sq"
)

var OrderBuilder = sq2.NewInsert().
	Into(sq2.NewTable("b2b.order")).
	SetColumns([]sq2.Column{
		sq2.NewColumn("id"),
		sq2.NewColumn("plan_id"),
		sq2.NewColumn("discount_id"),
		sq2.NewColumn("licence_id"),
		sq2.NewColumn("team_id"),
		sq2.NewColumn("checkout_id"),
		sq2.NewColumn("amount"),
		sq2.NewColumn("cycle_count"),
		sq2.NewColumn("trial_days"),
		sq2.NewColumn("kind"),
		sq2.NewColumn("created_utc"),
	})

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
