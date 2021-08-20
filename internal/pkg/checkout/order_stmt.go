package checkout

import "strings"

// StmtCreateBaseOrder is used to save an OrderRow
const StmtCreateBaseOrder = `
INSERT INTO b2b.order
SET id = :order_id,
	amount_payable = :amount_payable,
	created_by = :created_by,
	created_utc = :created_utc,
	item_count = :item_count,
	current_status = :current_status,
	team_id = :team_id,
	cart_items_summary = :cart_items_summary`

const StmtCreateOrderItem = `
INSERT INTO b2b.order_item
SET id = :order_item_id,
	order_id = :order_id,
	price_off_per_copy = :price_off_per_copy,
	price_snapshot = :price_snapshot,
	new_copy_count = :new_copy_count,
	renewal_licences = :renewal_licences`

const colBaseOrder = `
SELECT o.id AS order_id,
	o.amount_payable,
	o.created_by,
	o.created_utc,
	o.current_status,
	o.item_count,
	o.current_status,
	o.team_id
`

const colOrderList = colBaseOrder + `,
o.cart_items_summary`

const colOrderTeam = `
JSON_OBJECT(
	"orgName", t.org_name,
	"invoiceTitle", t.invoice_title,
	"phone", t.phone
)
AS team
`

//const StmtListOrders = colOrderList + `
//FROM b2b.order AS o
//WHERE o.team_id = ?
//ORDER BY o.created_utc DESC
//LIMIT ? OFFSET ?`

// BuildStmtListOrders retrieve multiple rows of order
// for CMS.
func BuildStmtListOrders(where string) string {
	return colOrderList + `,
` + colOrderTeam + `
FROM b2b.order AS o
LEFT JOIN b2b.team AS t
	ON o.team_id = t.id
WHERE ` + where + `
ORDER BY o.created_utc DESC
LIMIT ? OFFSET ?
`
}

func BuildStmtCountOrder(where string) string {
	return `
SELECT COUNT(*)
FROM b2b.order AS o
WHERE ` + where
}

const colPayment = `
o.amount_paid,
o.approved_by,
o.approved_utc,
o.description,
o.payment_method,
o.transaction_id
`

// BuildStmtOrder retrieve a row from order table
// without the cart_items_summary column.
// The data consists of the Payment and BaseOrder parts of
// an Order.
// withTeam - if true, adn team_id to where clause when used
// by an admin. When used by CMS, team_id should not present.
func BuildStmtOrder(withTeam bool) string {
	var buf strings.Builder

	buf.WriteString(colBaseOrder)
	buf.WriteByte(',')
	buf.WriteString(colPayment)
	buf.WriteString("FROM b2b.order AS o WHERE o.id = ?")
	if withTeam {
		buf.WriteString(" AND o.team_id = ?")
	}

	buf.WriteString(" LIMIT 1")

	return buf.String()
}

// StmtItemsOfOrder retrieves all items belong to an Order.
const StmtItemsOfOrder = `
SELECT id AS order_item_id,
	order_id,
	price_off_per_copy,
	price_snapshot,
	new_copy_count,
	renewal_licences
FROM b2b.order_item
WHERE order_id = ?
ORDER BY tier ASC`
