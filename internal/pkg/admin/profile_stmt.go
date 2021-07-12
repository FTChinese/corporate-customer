package admin

const StmtProfile = `
SELECT a.id AS admin_id,
	a.email AS email,
	a.display_name AS display_name,
	a.is_active AS active,
	a.verified AS verified,
	t.id AS team_id,
	t.org_name,
	t.invoice_tile,
FROM b2b.admin AS a
	LEFT JOIN b2b.team AS t
	ON a.id = t.admin_id
WHERE a.id = ?
LIMIT 1`
