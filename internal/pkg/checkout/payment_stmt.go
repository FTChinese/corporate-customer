package checkout

const StmtSavePayment = `
INSERT INTO b2b.payment
SET order_id = :order_id,
	amount_paid = :amount_paid,
	approved_by = :approved_by,
	approved_utc = :approved_utc,
	description = :description,
	payment_method = :payment_method,
	transaction_id = :transaction_id
`

const StmtSavePaymentOffer = `
INSERT INTO b2b.payment_offer
SET order_id = :order_id,
	copy_count = :copy_count,
	kind = :kind,
	price = :price,
	price_off_per_copy = :price_off_per_copy
`

// StmtListTxnToConfirm retrieve a licence queue and the price
// to bulk process licence creation/renewal.
const StmtListTxnToConfirm = `
SELECT t.transaction_id AS txn_id,
	i.price AS price
FROM b2b.licence_transaction AS t
	LEFT JOIN b2b.order_item AS i
ON t.order_id = i.order_id
	AND t.price_id = i.price_id	
WHERE t.order_id = ?
`

const StmtInsertPaymentErr = `
INSERT INTO b2b.payment_error
SET transaction_id = :txn_id,
	error_message = :error_message,
	created_utc = :created_utc
`
