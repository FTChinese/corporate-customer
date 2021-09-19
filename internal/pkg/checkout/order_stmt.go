package checkout

import "strings"

// StmtCreateOrder is used to save a row of Order.
const StmtCreateOrder = `
INSERT INTO b2b.order
SET id = :order_id,
	admin_id = :admin_id,
	team_id = :team_id,
	amount_payable = :amount_payable,
	created_utc = :created_utc,
	current_status = :current_status,
	item_count = :item_count,
	item_list = :item_list
`

// Shared columns used both when retrieving a list orders,
// or a single row of order.
const colOrder = `
SELECT o.id AS order_id,
	o.admin_id AS admin_id,
	o.team_id AS team_id,
	o.amount_payable AS amount_payable,
	o.created_utc AS created_utc,
	o.item_count AS item_count,
	o.item_list AS item_list,
	o.current_status AS current_status
`

// BuildStmtOrder retrieve a row from order table
// withTeam - if true, add team_id to where clause when used
// by an admin.
// When used by CMS, team_id should be omitted.
func BuildStmtOrder(withTeam bool) string {
	var buf strings.Builder

	buf.WriteString(colOrder)
	buf.WriteString("FROM b2b.order AS o WHERE o.id = ?")
	if withTeam {
		buf.WriteString(" AND o.team_id = ?")
	}

	buf.WriteString(" LIMIT 1")

	return buf.String()
}

// StmtListOrders retrieves a list of orders for an admin.
const StmtListOrders = colOrder + `
FROM b2b.order AS o
WHERE o.team_id = ?
ORDER BY o.created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountOrder = `
SELECT COUNT(*)
FROM b2b.order AS o
WHERE o.team_id = ?
`

// BuildStmtListOrdersCMS retrieve multiple rows of order
// for CMS.
// The WHERE clause is empty by default.
// Client could provide team_id or current_status or both
// to filter data.
func BuildStmtListOrdersCMS(where string) string {
	return colOrder + `,
JSON_OBJECT(
	"orgName", t.org_name,
	"invoiceTitle", t.invoice_title,
	"phone", t.phone
) AS team
FROM b2b.order AS o
LEFT JOIN b2b.team AS t
	ON o.team_id = t.id
` + where + `
ORDER BY o.created_utc DESC
LIMIT ? OFFSET ?
`
}

func BuildStmtCountOrder(where string) string {
	return `
SELECT COUNT(*)
FROM b2b.order AS o
` + where
}

const StmtUpdateOrderStatus = `
UPDATE b2b.order
SET current_status = :current_status
WHERE id = :order_id
LIMIT 1
`
