package order

const StmtCreateCheckoutOrder = `
INSERT INTO b2b.order
SET id = :order_id,
	amount_payable = :amount_payable,
	created_by = :created_by,
	checkout_products = :checkout_products,
	current_status = :current_status,
	item_count = :item_count,
	team_id = :team_id`

const StmtListCheckoutOrders = `
SELECT id AS order_id,
	amount_payable,
	created_by,
	created_utc,
	checkout_products,
	current_status,
	item_count,
	team_id
FROM b2b.order
WHERE AND team_id = ?
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtCreateOrderItem = `
INSERT INTO b2b.order_item
SET id = :order_item_id,
	order_id = :order_id,
	price_off_per_copy = :price_off_per_copy,
	price_snapshot = :price_snapshot,
	new_copy_count = :new_copy_count,
	renewal_licences = :renewal_licences`

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

// Deprecated
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

// Deprecated
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
