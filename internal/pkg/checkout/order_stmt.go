package checkout

const StmtCreateBaseOrder = `
INSERT INTO b2b.order
SET id = :order_id,
	amount_payable = :amount_payable,
	created_by = :created_by,
	checkout_products = :checkout_products,
	current_status = :current_status,
	item_count = :item_count,
	team_id = :team_id`

const StmtCreateOrderItem = `
INSERT INTO b2b.order_item
SET id = :order_item_id,
	order_id = :order_id,
	price_off_per_copy = :price_off_per_copy,
	price_snapshot = :price_snapshot,
	new_copy_count = :new_copy_count,
	renewal_licences = :renewal_licences`

const colOrder = `
SELECT id AS order_id,
	amount_payable,
	created_by,
	created_utc,
	current_status,
	item_count,
	team_id
`
const StmtListBaseOrders = colOrder + `,
	checkout_products
FROM b2b.order
WHERE team_id = ?
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountOrder = `
SELECT COUNT(*)
FROM b2b.order
WHERE team_id = ?`

const StmtOrderDetails = colOrder + `,
	amount_paid,
	approved_by,
	approved_utc,
	description,
	payment_method,
	transaction_id
FROM b2b.order
WHERE id = ?
	AND team_id = ?
LIMIT 1`

const StmtItemsOfOrder = `
SELECT order_item_id,
	order_id,
	price_off_per_copy,
	price_snapshot,
	new_copy_count,
	renewal_licences
FROM b2b.order_item
WHERE order_id = ?
ORDER BY tier ASC`
