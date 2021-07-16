package stmt

import (
	"github.com/FTChinese/ftacademy/pkg/sq"
)

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
