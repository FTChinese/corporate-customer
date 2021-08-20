package admin

import "strings"

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
FROM b2b.team
`

// BuildStmtLoadTeam build the sql statement to load a team.
// if adminOnly is true, this intends to be used by a loggedin
// admin only:
// WHERE id = ? AND admin_id = ? LIMIT 1`
//
// For CMS you should not add the admin_id condition in WHERE:
// WHERE id = ? LIMIT 1`
func BuildStmtLoadTeam(adminOnly bool) string {
	var buf strings.Builder

	buf.WriteString(colTeam)
	buf.WriteString("WHERE id = ?")
	if adminOnly {
		buf.WriteString(" AND admin_id = ?")
	}
	buf.WriteString(" LIMIT 1")

	return buf.String()
}

const StmtUpdateTeam = `
UPDATE b2b.team
SET org_name = :org_name,
	invoice_title = :invoice_title,
	phone = :phone,
	updated_utc = :updated_utc
WHERE id = :team_id
	AND admin_id = :admin_id
LIMIT 1`
