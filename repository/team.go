package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
	gorest "github.com/FTChinese/go-rest"
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

func (env Env) SaveTeamMember(m admin.TeamMember) error {
	_, err := env.db.NamedExec(stmt.InsertTeamMember, m)

	if err != nil {
		return err
	}

	return nil
}

// UpdateTeamMember add a member's ftc if missing.
// This is used after a reader signup upon verifying invitation.
func (env Env) UpdateTeamMember(m admin.TeamMember) error {
	_, err := env.db.NamedExec(stmt.SetTeamMemberFtcID, m)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ListTeamMembers(teamID string, page gorest.Pagination) ([]admin.TeamMember, error) {
	list := make([]admin.TeamMember, 0)

	err := env.db.Select(&list, stmt.ListTeamMembers, teamID, page.Limit, page.Offset())
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (env Env) AsyncListTeamMembers(teamID string, page gorest.Pagination) <-chan admin.TeamMemberList {
	r := make(chan admin.TeamMemberList)

	go func() {
		defer close(r)

		list, err := env.ListTeamMembers(teamID, page)

		r <- admin.TeamMemberList{
			Data: list,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountTeamMembers(teamID string) (int64, error) {
	var total int64
	err := env.db.Get(&total, stmt.CountTeamMembers, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountTeamMembers(teamID string) <-chan admin.TeamMemberList {
	r := make(chan admin.TeamMemberList)

	go func() {
		defer close(r)
		total, err := env.CountTeamMembers(teamID)

		r <- admin.TeamMemberList{
			Total: total,
			Err:   err,
		}
	}()

	return r
}

func (env Env) DeleteTeamMember(m admin.TeamMember) error {
	_, err := env.db.NamedExec(stmt.DeleteTeamMember, m)
	if err != nil {
		return err
	}

	return nil
}
