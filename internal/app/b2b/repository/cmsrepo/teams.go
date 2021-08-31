package cmsrepo

import "github.com/FTChinese/ftacademy/internal/pkg/admin"

func (env Env) LoadTeam(teamID string) (admin.Team, error) {
	var t admin.Team
	err := env.DBs.Read.Get(
		&t,
		admin.BuildStmtLoadTeam(false),
		teamID)
	if err != nil {
		return admin.Team{}, err
	}

	return t, nil
}
