package subs

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	gorest "github.com/FTChinese/go-rest"
)

// LoadLicence retrieves a licence, together with its
// subscription plan and the user to whom it was assigned.
// If the licence is not assigned yet, assignee fields are empty..
func (env Env) LoadLicence(r admin.AccessRight) (licence.Licence, error) {
	var lic licence.Licence
	err := env.dbs.Read.Get(&lic, licence.StmtLicence, r.RowID, r.TeamID)

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

func (env Env) ListLicence(teamID string, page gorest.Pagination) (licence.LicList, error) {
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
			PagedList: pkg.PagedList{
				Total:      0,
				Pagination: gorest.Pagination{},
				Err:        err,
			},
			Data: licences,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return licence.LicList{}, listResult.Err
	}
	return licence.LicList{
		PagedList: pkg.PagedList{
			Total:      count,
			Pagination: page,
			Err:        nil,
		},
		Data: listResult.Data,
	}, nil
}

// GrantLicence grants a licence to a reader.
func (env Env) GrantLicence(r admin.AccessRight, to licence.Assignee) (licence.GrantResult, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.GrantResult{}, err
	}
	// Retrieve the licence to be granted.
	lic, err := tx.RetrieveBaseLicence(r)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}
	if !lic.IsAvailable() {
		_ = tx.Rollback()
		return licence.GrantResult{}, errors.New("licence is not available to grant")
	}

	// Retrieve the invitation of this session.
	inv, err := tx.RetrieveInvitation(admin.AccessRight{
		RowID:  lic.LatestInvitation.ID,
		TeamID: r.TeamID,
	})
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}
	if !inv.IsAcceptable() {
		_ = tx.Rollback()
		return licence.GrantResult{}, errors.New("invitation is not acceptable")
	}

	// Retrieve membership for this user.
	mmb, err := tx.RetrieveMember(to.FtcID.String)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// TODO: Ensure membership could use a licence.

	// Update invitation.
	acceptedInv := inv.Accepted()
	// Update licence.
	grantedLic := lic.Granted(to, acceptedInv)
	// Create/update membership based on licence.
	result := licence.NewGrantResult(licence.Licence{
		BaseLicence: grantedLic,
		Assignee:    licence.AssigneeJSON{Assignee: to},
	}, mmb)

	// Update invitation
	err = tx.UpdateInvitationStatus(acceptedInv)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// Update licence
	err = tx.UpdateLicenceStatus(grantedLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// Insert a new row if the current membership is empty.
	if !mmb.IsZero() {
		err := tx.CreateMember(result.Membership)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return licence.GrantResult{}, err
		}
	} else {
		// Update current membership based on
		// licence.
		err := tx.UpdateMember(result.Membership)
		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return licence.GrantResult{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return licence.GrantResult{}, err
	}

	// The returned data is used to compose a letter
	return result, nil
}

// RevokeLicence revokes a licence granted to a reader.
// For a licence waiting for its invitation accepted,
// use RevokeInvitation instead of this one.
// After licence revoked, send email both to admin and
// ex-owner to notify this event.
func (env Env) RevokeLicence(r admin.AccessRight) (licence.RevokeResult, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.RevokeResult{}, err
	}

	lic, err := tx.RetrieveBaseLicence(r)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, err
	}
	if lic.IsRevocable() {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, errors.New("nothing to revoke")
	}

	mmb, err := tx.RetrieveMember(lic.AssigneeID.String)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, err
	}
	if !lic.IsGrantedTo(mmb) {
		_ = tx.Rollback()
		return licence.RevokeResult{}, errors.New("reader's membership is not generated from b2b licence")
	}

	updatedLic := lic.Revoked()
	updatedMmb := licence.RevokeLicence(mmb)

	err = tx.UpdateLicenceStatus(updatedLic)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, nil
	}

	err = tx.UpdateMember(updatedMmb)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, nil
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return licence.RevokeResult{}, err
	}

	return licence.RevokeResult{
		Licence: licence.Licence{
			BaseLicence: updatedLic,
			Assignee:    licence.AssigneeJSON{},
		},
		Membership: updatedMmb,
		Snapshot:   mmb.Snapshot(reader.B2BArchiver(reader.ArchiveActionRevoke)),
	}, nil
}
