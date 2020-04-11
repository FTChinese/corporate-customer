package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// CreateTeam creates a new team under an admin account.
// We need to make sure the admin actually exists before
// creating a team.
func (env Env) CreateTeam(t admin.Team) error {
	_, err := env.db.NamedExec(stmt.CreateTeam, t)

	if err != nil {
		return err
	}

	return nil
}

// TeamById retrieves a team by its id.
func (env Env) TeamByAdminID(adminID string) (admin.Team, error) {
	var t admin.Team

	err := env.db.Get(&t, stmt.TeamByAdminID, adminID)

	if err != nil {
		return t, err
	}

	return t, nil
}

// UpdateTeam changes a team's name to invoice title.
func (env Env) UpdateTeam(t admin.Team) error {
	_, err := env.db.NamedExec(stmt.UpdateTeam, t)

	if err != nil {
		return err
	}

	return nil
}
