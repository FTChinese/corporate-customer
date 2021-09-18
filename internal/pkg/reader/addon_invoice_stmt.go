package reader

const StmtCreateInvoice = `
INSERT INTO premium.ftc_invoice
SET id = :id,
	user_compound_id = :compound_id,
	tier = :tier,
	cycle = :cycle,
	years = :years,
	months = :months,
	extra_days = :days,
	addon_source = :addon_source,
	apple_tx_id = :apple_tx_id,
	licence_tx_id = :licence_tx_id,
	order_id = :order_id,
	order_kind = :order_kind,
	paid_amount = :paid_amount,
	payment_method = :payment_method,
	price_id = :price_id,
	stripe_subs_id = :stripe_subs_id,
	created_utc = :created_utc,
	consumed_utc = :consumed_utc,
	start_utc = :start_utc,
	end_utc = :end_utc,
	carried_over_utc = :carried_over_utc
`
