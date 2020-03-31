package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

const stmtCreateTeam = `
INSERT INTO b2b.team
SET id = :team_id,
	name = :name,
	invoice_title = :invoice_title,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP(),
	admin_id = :admin_id`

// CreateTeam creates a new team under an admin account.
// We need to make sure the admin actually exists before
// creating a team.
func (env Env) CreateTeam(t admin.Team) error {
	_, err := env.db.NamedExec(stmtCreateTeam, t)

	if err != nil {
		return err
	}

	return nil
}

const stmtTeamByID = stmt.TeamBase + `
WHERE id = ?
LIMIT 1`

// TeamById retrieves a team by its id.
func (env Env) TeamByID(teamID string) (admin.Team, error) {
	var t admin.Team

	err := env.db.Get(&t, stmtTeamByID, teamID)

	if err != nil {
		return t, err
	}

	return t, nil
}

const stmtTeamOfAdmin = stmt.TeamBase + `
WHERE admin_id = ?
LIMIT 1`

// TeamOfAdmin retrieves a team belong to an admin.
func (env Env) TeamOfAdmin(adminID string) (admin.Team, error) {
	var t admin.Team

	err := env.db.Get(&t, stmtTeamOfAdmin, adminID)

	if err != nil {
		return t, err
	}

	return t, nil
}

const stmtUpdateTeam = `
UPDATE b2b.team
SET name = :name,
	invoice_title = :invoice_title
WHERE id = ?
LIMIT 1`

// UpdateTeam changes a team's name to invoice title.
func (env Env) UpdateTeam(t admin.Team) error {
	_, err := env.db.NamedExec(stmtUpdateTeam, t)

	if err != nil {
		return err
	}

	return nil
}
