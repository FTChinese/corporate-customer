package checkout

const StmtInsertCartItem = `
INSERT INTO b2b.order_item
SET order_id = :order_id,
	price = :price,
	new_copy_count = :new_copy_count,
	renewal_list = :renewal_list,
	admin_id = :admin_id,
	team_id = :team_id,
	created_utc = :created_utc
`
