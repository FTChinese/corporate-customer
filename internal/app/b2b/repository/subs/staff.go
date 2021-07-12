package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) SaveStaffer(m model2.Staffer) error {
	_, err := env.dbs.Write.NamedExec(stmt.InsertStaffer, m)

	if err != nil {
		return err
	}

	return nil
}

// UpdateStaffer add a member's ftc if missing.
// This is used after a reader signup upon verifying invitation.
func (env Env) UpdateStaffer(m model2.Staffer) error {
	_, err := env.dbs.Write.NamedExec(stmt.SetStaffFtcID, m)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ListStaff(teamID string, page gorest.Pagination) ([]model2.Staffer, error) {
	list := make([]model2.Staffer, 0)

	err := env.dbs.Read.Select(&list, stmt.ListStaff, teamID, page.Limit, page.Offset())
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (env Env) AsyncListStaff(teamID string, page gorest.Pagination) <-chan model2.StaffList {
	r := make(chan model2.StaffList)

	go func() {
		defer close(r)

		list, err := env.ListStaff(teamID, page)

		r <- model2.StaffList{
			Data: list,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountStaff(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, stmt.CountStaff, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountStaff(teamID string) <-chan model2.StaffList {
	r := make(chan model2.StaffList)

	go func() {
		defer close(r)
		total, err := env.CountStaff(teamID)

		r <- model2.StaffList{
			Total: total,
			Err:   err,
		}
	}()

	return r
}

// DeleteStaffer deletes a staffer that is not a member of a team.
// TODO: A this staffer is still using a licence of this team,
// delete should be ignored.
func (env Env) DeleteStaffer(m model2.Staffer) error {
	_, err := env.dbs.Write.NamedExec(stmt.DeleteStaffer, m)
	if err != nil {
		return err
	}

	return nil
}
