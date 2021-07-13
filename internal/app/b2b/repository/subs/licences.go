package subs

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	gorest "github.com/FTChinese/go-rest"
)

// LoadLicence retrieves a licence, together with its
// subscription plan and the user to whom it was assigned.
// If the licence is not assigned yet, assignee fields are empty..
func (env Env) LoadLicence(id, teamID string) (licence.Licence, error) {
	var lic licence.Licence
	err := env.dbs.Read.Get(&lic, licence.StmtLicence, id, teamID)

	if err != nil {
		return licence.Licence{}, err
	}

	return lic, nil
}

func (env Env) LicenceOfInvitation(id string) (licence.Licence, error) {
	var lic licence.Licence
	err := env.dbs.Read.Get(&lic, licence.StmtInvitedLicence, id)
	if err != nil {
		return licence.Licence{}, err
	}

	return lic, nil
}

// listLicences shows a list all licence.
// Each licence's plan, invitation, assignee are attached.
func (env Env) listLicences(teamID string, page gorest.Pagination) ([]licence.Licence, error) {
	var licences = make([]licence.Licence, 0)

	err := env.dbs.Read.Select(
		&licences,
		licence.StmtListLicences,
		teamID,
		page.Limit,
		page.Offset(),
	)

	if err != nil {
		return nil, err
	}

	return licences, nil
}

func (env Env) countLicences(teamID string) (int64, error) {
	var total int64
	if err := env.dbs.Read.Get(&total, licence.StmtCountLicence, teamID); err != nil {
		return total, err
	}

	return total, nil
}

func (env Env) ListExpLicence(teamID string, page gorest.Pagination) (licence.LicList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan licence.LicList)

	go func() {
		defer close(countCh)
		n, err := env.countLicences(teamID)
		if err != nil {
			sugar.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		licences, err := env.listLicences(teamID, page)

		listCh <- licence.LicList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       licences,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return licence.LicList{}, listResult.Err
	}
	return licence.LicList{
		Total:      count,
		Pagination: page,
		Data:       listResult.Data,
	}, nil
}
