package adminor

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
)

// CreateTeam creates a new team under an admin account.
// We need to make sure the admin actually exists before
// creating a team.
func (env Env) CreateTeam(t admin.Team) error {
	_, err := env.DBs.Write.NamedExec(admin.StmtCreateTeam, t)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadTeam(id string) (admin.Team, error) {
	var t admin.Team
	err := env.DBs.Read.Get(&t, admin.StmtTeamByID, id)
	if err != nil {
		return admin.Team{}, err
	}

	return t, nil
}

// UpdateTeam changes a team's name to invoice title.
func (env Env) UpdateTeam(t admin.Team) error {
	_, err := env.DBs.Write.NamedExec(admin.StmtUpdateTeam, t)

	if err != nil {
		return err
	}

	return nil
}
