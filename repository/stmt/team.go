package stmt

const InsertTeam = `
INSERT INTO b2b.team
SET id = :team_id,
	name = :name,
	invoice_title = :invoice_title,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP(),
	admin_id = :admin_id`

const RetrieveTeam = `
SELECT id AS team_id
	name,
	invoice_tile,
	created_utc
FROM b2b.team
WHERE admin_id = ?
LIMIT 1`

const UpdateTeam = `
UPDATE b2b.team
SET name = :name
	display_name = :display_name
WHERE id = ?
LIMIT 1`
