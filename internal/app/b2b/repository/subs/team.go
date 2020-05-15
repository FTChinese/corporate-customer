package subs

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/b2b/internal/app/b2b/stmt"
	gorest "github.com/FTChinese/go-rest"
)

// CreateTeam creates a new team under an admin account.
// We need to make sure the admin actually exists before
// creating a team.
func (env Env) CreateTeam(t model.Team) error {
	_, err := env.db.NamedExec(stmt.CreateTeam, t)

	if err != nil {
		return err
	}

	return nil
}

// TeamById retrieves a team by its id.
func (env Env) TeamByAdminID(adminID string) (model.Team, error) {
	var t model.Team

	err := env.db.Get(&t, stmt.TeamByAdminID, adminID)

	if err != nil {
		return t, err
	}

	return t, nil
}

// UpdateTeam changes a team's name to invoice title.
func (env Env) UpdateTeam(t model.Team) error {
	_, err := env.db.NamedExec(stmt.UpdateTeam, t)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) SaveStaffer(m model.Staffer) error {
	_, err := env.db.NamedExec(stmt.InsertStaffer, m)

	if err != nil {
		return err
	}

	return nil
}

// UpdateStaffer add a member's ftc if missing.
// This is used after a reader signup upon verifying invitation.
func (env Env) UpdateStaffer(m model.Staffer) error {
	_, err := env.db.NamedExec(stmt.SetStaffFtcID, m)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ListStaff(teamID string, page gorest.Pagination) ([]model.Staffer, error) {
	list := make([]model.Staffer, 0)

	err := env.db.Select(&list, stmt.ListStaff, teamID, page.Limit, page.Offset())
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (env Env) AsyncListStaff(teamID string, page gorest.Pagination) <-chan model.StaffList {
	r := make(chan model.StaffList)

	go func() {
		defer close(r)

		list, err := env.ListStaff(teamID, page)

		r <- model.StaffList{
			Data: list,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountStaff(teamID string) (int64, error) {
	var total int64
	err := env.db.Get(&total, stmt.CountStaff, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountStaff(teamID string) <-chan model.StaffList {
	r := make(chan model.StaffList)

	go func() {
		defer close(r)
		total, err := env.CountStaff(teamID)

		r <- model.StaffList{
			Total: total,
			Err:   err,
		}
	}()

	return r
}

// DeleteStaffer deletes a staffer that is not a member of a team.
// TODO: A this staffer is still using a licence of this team,
// delete should be ignored.
func (env Env) DeleteStaffer(m model.Staffer) error {
	_, err := env.db.NamedExec(stmt.DeleteStaffer, m)
	if err != nil {
		return err
	}

	return nil
}
