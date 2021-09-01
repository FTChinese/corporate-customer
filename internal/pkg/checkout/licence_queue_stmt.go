package checkout

import "github.com/FTChinese/ftacademy/pkg/sq"

// StmtBulkLicenceQueue build the SQL to bulk insert into
// licence_queue table.
// Its value is built by BulkLicenceQueue.
func StmtBulkLicenceQueue(n int) sq.BulkInsert {
	return sq.NewBulkInsert().
		Into(sq.NewTable("b2b.licence_queue")).
		SetColumns(
			sq.NewColumn("finalized_utc"),
			sq.NewColumn("array_index"),
			sq.NewColumn("kind"),
			sq.NewColumn("licence_prior"),
			sq.NewColumn("licence_after"),
			sq.NewColumn("order_id"),
			sq.NewColumn("price_id"),
			sq.NewColumn("created_utc"),
		).
		Rows(n)
}

const StmtFinalizeLicenceQueue = `
UPDATE b2b.licence_queue
SET finalized_utc = :finalized_utc,
	licence_after = :licence_after
WHERE id = :id
LIMIT 1`

const colLicenceQueue = `
SELECT id,
	created_utc,
	finalized_utc,
	array_index,
	kind,
	licence_prior,
	licence_after,
	order_id,
	price_id
FROM b2b.licence_queue
`

// StmtListLicenceQueue retrieves all queued licence of
// a specific price of an order.
const StmtListLicenceQueue = colLicenceQueue + `
WHERE order_id = ?
	AND price_id = ?
ORDER BY kind ASC, array_index ASC
`

// StmtLockLicenceQueue locks a row in licence queue
// when payment is confirmed.
const StmtLockLicenceQueue = colLicenceQueue + `
WHERE id = ?
LIMIT 1
FOR UPDATE
`

// StmtListPriceOfQueue retrieve a licence queue and the price
// to bulk process licence creation/renewal.
const StmtListPriceOfQueue = `
SELECT q.id AS queue_id,
	c.price AS price
FROM b2b.licence_queue AS q
	LEFT JOIN cart_item AS c
ON q.order_id = c.order_id
	AND q.price_id = c.price_id	
WHERE q.order_id = ?
`
