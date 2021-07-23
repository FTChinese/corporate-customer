package admin

const StmtCreateTeam = `
INSERT INTO b2b.team
SET id = :team_id,
	admin_id = :admin_id,
	org_name = :org_name,
	invoice_title = :invoice_title,
	created_utc = :created_utc`

const colTeam = `
SELECT id AS team_id,
	admin_id,
	org_name,
	invoice_title,
	created_utc
FROM b2b.team`

const StmtTeamByAdminID = colTeam + `
WHERE admin_id = ?
LIMIT 1`

const StmtTeamByID = colTeam + `
WHERE id = ?
LIMIT 1`

const StmtUpdateTeam = `
UPDATE b2b.team
SET name = :name,
	invoice_title = :invoice_title
WHERE id = ?
LIMIT 1`
