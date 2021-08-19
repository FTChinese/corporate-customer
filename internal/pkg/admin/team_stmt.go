package admin

const StmtCreateTeam = `
INSERT INTO b2b.team
SET id = :team_id,
	admin_id = :admin_id,
	org_name = :org_name,
	phone = :phone,
	invoice_title = :invoice_title,
	created_utc = :created_utc`

const colTeam = `
SELECT id AS team_id,
	admin_id,
	org_name,
	phone,
	invoice_title,
	created_utc
FROM b2b.team`

// StmtTeamOfAdmin retrieve a row by id belong to an admin.
// Used by ajax api.
const StmtTeamOfAdmin = colTeam + `
WHERE id = ?
	AND admin_id = ?
LIMIT 1`

// StmtTeamByID retrieve a row by id.
// Used by restful api.
const StmtTeamByID = colTeam + `
WHERE id = ?
LIMIT 1`

const StmtUpdateTeam = `
UPDATE b2b.team
SET org_name = :org_name,
	invoice_title = :invoice_title,
	phone = :phone,
	updated_utc = :updated_utc
WHERE id = :team_id
	AND admin_id = :admin_id
LIMIT 1`
