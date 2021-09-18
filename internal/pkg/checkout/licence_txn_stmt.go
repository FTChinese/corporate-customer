package checkout

import "github.com/FTChinese/ftacademy/pkg/sq"

// StmtBulkLicenceTxn build the SQL to bulk insert into
// licence_queue table.
// Its value is built by BulkLicenceTxn.
func StmtBulkLicenceTxn(n int) sq.BulkInsert {
	return sq.NewBulkInsert().
		Into(sq.NewTable("b2b.licence_transaction")).
		SetColumns(
			sq.NewColumn("transaction_id"),
			sq.NewColumn("kind"),
			sq.NewColumn("licence_to_renew"),
			sq.NewColumn("order_id"),
			sq.NewColumn("price_id"),
			sq.NewColumn("admin_id"),
			sq.NewColumn("team_id"),
			sq.NewColumn("created_utc"),
			sq.NewColumn("finalized_utc"),
		).
		Rows(n)
}

const colLicenceTxn = `
SELECT transaction_id AS txn_id,
	kind,
	licence_to_renew,
	order_id,
	price_id,
	admin_id,
	team_id,
	created_utc,
	finalized_utc
FROM b2b.licence_transaction
`

// StmtListLicenceTxn retrieves all queued licence of
// a specific price of an order.
const StmtListLicenceTxn = colLicenceTxn + `
WHERE order_id = ?
	AND price_id = ?
ORDER BY id ASC
`

// StmtLockLicenceTxn locks a row in licence queue
// when payment is confirmed.
const StmtLockLicenceTxn = colLicenceTxn + `
WHERE transaction_id = ?
LIMIT 1
FOR UPDATE
`

const StmtFinalizeLicenceTxn = `
UPDATE b2b.licence_transaction
SET finalized_utc = :finalized_utc
WHERE transaction_id = :txn_id
LIMIT 1
`

const colLicenceUpsert = `
current_period_start_utc = :current_period_start_utc,
current_period_end_utc = :current_period_end_utc,
hint_grant_mismatch = :hint_grant_mismatch,
latest_transaction_id = :latest_transaction_id,
latest_price = :latest_price,
updated_utc = :updated_utc`

const StmtCreateLicence = `
INSERT INTO b2b.licence
SET id = :licence_id,
	tier = :tier,
	cycle = :cycle,
	current_status = :lic_status,
	admin_id = :admin_id,
	team_id = :team_id,

	start_date_utc = :start_date_utc,
	trial_start_utc = :trial_start_utc,
	trial_end_utc = :trial_end_utc,
	
	latest_invitation = :latest_invitation,
	assignee_id = :assignee_id,
	created_utc = :created_utc,
` + colLicenceUpsert

const StmtRenewLicence = `
UPDATE b2b.licence
SET 
` + colLicenceUpsert + `
WHERE id = :licence_id
LIMIT 1`
