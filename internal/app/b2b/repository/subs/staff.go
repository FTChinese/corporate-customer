package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	gorest "github.com/FTChinese/go-rest"
)

func (env Env) SaveStaffer(m licence.Staffer) error {
	_, err := env.dbs.Write.NamedExec(licence.InsertStaffer, m)

	if err != nil {
		return err
	}

	return nil
}

// UpdateStaffer add a member's ftc if missing.
// This is used after a reader signup upon verifying invitation.
func (env Env) UpdateStaffer(m licence.Staffer) error {
	_, err := env.dbs.Write.NamedExec(licence.SetStaffFtcID, m)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ListStaff(teamID string, page gorest.Pagination) ([]licence.Staffer, error) {
	list := make([]licence.Staffer, 0)

	err := env.dbs.Read.Select(&list, licence.ListStaff, teamID, page.Limit, page.Offset())
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (env Env) AsyncListStaff(teamID string, page gorest.Pagination) <-chan licence.StaffList {
	r := make(chan licence.StaffList)

	go func() {
		defer close(r)

		list, err := env.ListStaff(teamID, page)

		r <- licence.StaffList{
			Data: list,
			Err:  err,
		}
	}()

	return r
}

func (env Env) CountStaff(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, licence.CountStaff, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) AsyncCountStaff(teamID string) <-chan licence.StaffList {
	r := make(chan licence.StaffList)

	go func() {
		defer close(r)
		total, err := env.CountStaff(teamID)

		r <- licence.StaffList{
			Total: total,
			Err:   err,
		}
	}()

	return r
}

// DeleteStaffer deletes a staffer that is not a member of a team.
// TODO: A this staffer is still using a licence of this team,
// delete should be ignored.
func (env Env) DeleteStaffer(m licence.Staffer) error {
	_, err := env.dbs.Write.NamedExec(licence.DeleteStaffer, m)
	if err != nil {
		return err
	}

	return nil
}
