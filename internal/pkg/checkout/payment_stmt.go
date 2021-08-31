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
	array_index = :array_index,
	copy_count = :copy_count,
	kind = :kind,
	price = :price,
	price_off_per_copy = :price_off_per_copy
`
