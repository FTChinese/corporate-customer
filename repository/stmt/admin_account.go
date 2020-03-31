// Package stmt provides shared, reusable SQL statements.
// Most of them cannot be executed alone.
package stmt

const AccountBase = `
SELECT a.id AS admin_id,
	a.email AS email,
	a.display_name AS display_name,
	a.is_active AS active,
	a.verified AS verified,
	a.created_utc AS created_utc,
	a.updated_utc AS updated_utc`

const TeamBase = `
SELECT id AS team_id
	name,
	invoice_tile,
	created_utc
FROM b2b.team`
