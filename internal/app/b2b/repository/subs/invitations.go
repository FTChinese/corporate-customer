package subs

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	gorest "github.com/FTChinese/go-rest"
)

// InvitationByToken tries to find an Invitation by token.
func (env Env) InvitationByToken(token string) (licence.Invitation, error) {
	var inv licence.Invitation
	err := env.dbs.Read.Get(&inv, licence.StmtInvitationByToken, token)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

func (env Env) InvitationByID(r admin.AccessRight) (licence.Invitation, error) {
	var inv licence.Invitation
	err := env.dbs.Read.Get(&inv, licence.StmtInvitationByID, r.RowID, r.TeamID)
	if err != nil {
		return licence.Invitation{}, err
	}

	return inv, nil
}

// CreateInvitation creates an invitation for a licence
// depending on the licence availability.
// The returned licence contains the newly created invitation instance.
func (env Env) CreateInvitation(params input.InvitationParams, adminID string) (licence.BaseLicence, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.BaseLicence{}, err
	}

	lic, err := tx.RetrieveBaseLicence(admin.AccessRight{
		RowID:  params.LicenceID,
		TeamID: params.TeamID,
	})
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.BaseLicence{}, err
	}
	if !lic.IsAvailable() {
		_ = tx.Rollback()
		return licence.BaseLicence{}, err
	}

	inv, err := licence.NewInvitation(params, adminID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.BaseLicence{}, err
	}

	updateLic := lic.WithInvitation(inv)

	err = tx.CreateInvitation(inv)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.BaseLicence{}, err
	}

	err = tx.UpdateLicenceStatus(updateLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.BaseLicence{}, err
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return licence.BaseLicence{}, err
	}

	return updateLic, nil
}

func (env Env) RevokeInvitation(invID, teamID string) (licence.Invitation, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.Invitation{}, err
	}

	inv, err := tx.RetrieveInvitation(admin.AccessRight{
		RowID:  invID,
		TeamID: teamID,
	})
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.Invitation{}, err
	}
	if !inv.IsRevocable() {
		_ = tx.Rollback()
		return licence.Invitation{}, errors.New("invitation is not revocable")
	}

	lic, err := tx.RetrieveBaseLicence(admin.AccessRight{
		RowID:  inv.LicenceID,
		TeamID: teamID,
	})
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.Invitation{}, err
	}
	if !lic.IsInvitationRevocable() {
		_ = tx.Rollback()
		return licence.Invitation{}, errors.New("invitation is not revocable")
	}

	revokedInv := inv.Revoked()
	updatedLic := lic.WithInvitationRevoked()

	err = tx.UpdateInvitationStatus(revokedInv)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.Invitation{}, err
	}

	err = tx.UpdateLicenceStatus(updatedLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.Invitation{}, err
	}

	if err := tx.Commit(); err != nil {
		return licence.Invitation{}, err
	}

	return revokedInv, nil
}

// List invitations shows a list of invitations for a team.
func (env Env) listInvitations(teamID string, page gorest.Pagination) ([]licence.Invitation, error) {
	var invs = make([]licence.Invitation, 0)

	err := env.dbs.Read.Select(
		&invs,
		licence.StmtListInvitation,
		teamID,
		page.Limit,
		page.Offset())

	if err != nil {
		return nil, err
	}

	return invs, nil
}

func (env Env) countInvitation(teamID string) (int64, error) {
	var total int64

	err := env.dbs.Read.Get(
		&total,
		licence.StmtCountInvitation,
		teamID)

	if err != nil {
		return total, err
	}

	return total, nil
}

func (env Env) ListInvitations(teamID string, page gorest.Pagination) (licence.InvitationList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan licence.InvitationList)

	go func() {
		defer close(countCh)
		n, err := env.countInvitation(teamID)
		if err != nil {
			sugar.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)

		invs, err := env.listInvitations(teamID, page)

		listCh <- licence.InvitationList{
			PagedList: pkg.PagedList{
				Total:      0,
				Pagination: gorest.Pagination{},
				Err:        err,
			},
			Data: invs,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return licence.InvitationList{}, listResult.Err
	}

	return licence.InvitationList{
		PagedList: pkg.PagedList{
			Total:      count,
			Pagination: page,
			Err:        nil,
		},
		Data: listResult.Data,
	}, nil
}
