package subsrepo

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
func (env Env) LoadLicence(r admin.AccessRight) (licence.ExpandedLicence, error) {
	var lic licence.ExpandedLicence
	err := env.DBs.Read.Get(&lic, licence.StmtLicence, r.RowID, r.TeamID)

	if err != nil {
		return licence.ExpandedLicence{}, err
	}

	return lic, nil
}

// listLicences shows a list all licence.
// Each licence's plan, invitation, assignee are attached.
func (env Env) listLicences(teamID string, page gorest.Pagination) ([]licence.ExpandedLicence, error) {
	var licences = make([]licence.ExpandedLicence, 0)

	err := env.DBs.Read.Select(
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
	if err := env.DBs.Read.Get(&total, licence.StmtCountLicence, teamID); err != nil {
		return total, err
	}

	return total, nil
}

func (env Env) ListLicence(teamID string, page gorest.Pagination) (licence.PagedLicenceList, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan licence.PagedLicenceList)

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

		listCh <- licence.PagedLicenceList{
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
		return licence.PagedLicenceList{}, listResult.Err
	}
	return licence.PagedLicenceList{
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
	lic, err := tx.LockLicence(r)
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
	mmb, err := tx.LockMember(to.FtcID.String)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// TODO: if
	result, err := licence.GrantLicence(licence.GrantParams{
		CurLic:    lic,
		CurInv:    inv,
		To:        to,
		CurMember: mmb,
	})

	if err != nil {
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// Update invitation
	err = tx.UpdateInvitationStatus(result.LicenceVersion.PostChange.LatestInvitation.Invitation)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// Update licence
	err = tx.UpdateLicenceStatus(result.LicenceVersion.PostChange.Licence)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	// Upsert a new row if the current membership is empty.
	// Is current membership is zero, do insert;
	// otherwise update.
	err = tx.UpsertMember(result.MembershipVersion.PostChange.Membership, mmb.IsZero())
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return licence.GrantResult{}, err
	}

	// The returned data is used to compose a letter
	return result, nil
}

func (env Env) RevokeReplacedLicence(mv reader.MembershipVersioned, teamID string) (licence.GrantResult, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.GrantResult{}, err
	}

	lic, err := tx.LockLicence(admin.AccessRight{
		RowID:  mv.AnteChange.B2BLicenceID.String,
		TeamID: teamID,
	})

	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.GrantResult{}, err
	}

	result := licence.GrantResult{
		LicenceVersion: lic.Revoked().
			Versioned(licence.VersionActionRevoke).
			WithPriorVersion(lic).
			WithMembershipVersioned(mv.ID),
		MemberModified: licence.MemberModified{},
	}

	err = tx.UpdateLicenceStatus(result.LicenceVersion.PostChange.Licence)
	if err != nil {
		return licence.GrantResult{}, err
	}

	_ = tx.Commit()

	return result, nil
}

// RevokeLicence revokes a licence granted to a reader.
// For a licence waiting for its invitation accepted,
// use RevokeInvitation instead of this one.
// After licence revoked, send email both to admin and
// ex-owner to notify this event.
// There 2 cases when revoking:
// * User is using the licence and not expired, clean the assignee field and set membership to a past time.
// * User is using the licence but expired, simply clean the assignee field.
// * User changed to other payment method, simply clean the assignee field.
func (env Env) RevokeLicence(r admin.AccessRight) (licence.RevokeResult, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.beginTx()
	if err != nil {
		sugar.Error(err)
		return licence.RevokeResult{}, err
	}

	lic, err := tx.LockLicence(r)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, err
	}
	sugar.Infof("ExpandedLicence to revoke: %v", lic)
	if !lic.IsRevocable() {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, errors.New("nothing to revoke")
	}

	mmb, err := tx.LockMember(lic.AssigneeID.String)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, err
	}
	sugar.Infof("AnteChange to revoke: %v", mmb)

	result, err := licence.RevokeLicence(lic, mmb)
	if err != nil {
		_ = tx.Rollback()
		return licence.RevokeResult{}, err
	}

	err = tx.UpdateLicenceStatus(result.LicenceVersion.PostChange.Licence)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, nil
	}

	err = tx.UpdateMember(result.MembershipVersioned.PostChange.Membership)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return licence.RevokeResult{}, nil
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return licence.RevokeResult{}, err
	}

	return result, nil
}
