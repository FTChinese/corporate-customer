package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
)

func (tx TxRepo) LockBaseLicence(r admin.AccessRight) (licence.Licence, error) {
	var bl licence.Licence
	err := tx.Get(
		&bl,
		licence.StmtLockBaseLicence,
		r.RowID,
		r.TeamID)

	if err != nil {
		return licence.Licence{}, err
	}

	return bl, nil
}

// UpdateLicenceStatus after its status changed.
// * when invitation is created;
// * when it is revoked
// * when it is granted
func (tx TxRepo) UpdateLicenceStatus(lic licence.Licence) error {
	_, err := tx.NamedExec(licence.StmtUpdateLicenceStatus, lic)
	if err != nil {
		return err
	}

	return nil
}

func (tx TxRepo) CreateInvitation(inv licence.Invitation) error {
	_, err := tx.NamedExec(licence.StmtCreateInvitation, inv)
	if err != nil {
		return err
	}

	return nil
}

func (tx TxRepo) RetrieveInvitation(r admin.AccessRight) (licence.Invitation, error) {
	var inv licence.Invitation
	err := tx.Get(
		&inv,
		licence.StmtLockInvitation,
		r.RowID,
		r.TeamID)
	if err != nil {
		return licence.Invitation{}, err
	}

	return inv, nil
}

func (tx TxRepo) UpdateInvitationStatus(inv licence.Invitation) error {
	_, err := tx.NamedExec(licence.StmtUpdateInvitationStatus, inv)
	if err != nil {
		return err
	}

	return nil
}
