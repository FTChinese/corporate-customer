package subsrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg"
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

func (env Env) listStaff(teamID string, page gorest.Pagination) ([]licence.Staffer, error) {
	list := make([]licence.Staffer, 0)

	err := env.dbs.Read.Select(&list, licence.ListStaff, teamID, page.Limit, page.Offset())
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (env Env) countStaff(teamID string) (int64, error) {
	var total int64
	err := env.dbs.Read.Get(&total, licence.CountStaff, teamID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (env Env) ListStaff(teamID string, page gorest.Pagination) (licence.StaffList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan licence.StaffList)

	go func() {
		defer close(countCh)

		n, err := env.countStaff(teamID)
		if err != nil {
			sugar.Error(err)
		}
		countCh <- n
	}()

	go func() {
		defer close(listCh)

		list, err := env.listStaff(teamID, page)

		listCh <- licence.StaffList{
			PagedList: pkg.PagedList{
				Total:      0,
				Pagination: gorest.Pagination{},
				Err:        err,
			},
			Data: list,
		}
	}()

	count, listResult := <-countCh, <-listCh
	if listResult.Err != nil {
		return licence.StaffList{}, listResult.Err
	}

	return licence.StaffList{
		PagedList: pkg.PagedList{
			Total:      count,
			Pagination: page,
			Err:        nil,
		},
		Data: listResult.Data,
	}, nil
}

// DeleteStaffer deletes a staffer that is not a member of a team.
func (env Env) DeleteStaffer(m licence.Staffer) error {
	_, err := env.dbs.Write.NamedExec(licence.DeleteStaffer, m)
	if err != nil {
		return err
	}

	return nil
}
