package stmt

const CreateTeam = `
INSERT INTO b2b.team
SET id = :team_id,
	admin_id = :admin_id,
	name = :name,
	invoice_title = :invoice_title,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const selectTeam = `
SELECT id AS team_id
	admin_id
	name,
	invoice_tile,
	created_utc
FROM b2b.team`

const TeamByAdminID = selectTeam + `
WHERE admin_id = ?
LIMIT 1`

const TeamByID = selectTeam + `
WHERE id = ?
LIMIT 1`

const UpdateTeam = `
UPDATE b2b.team
SET name = :name,
	invoice_title = :invoice_title
WHERE id = ?
LIMIT 1`
