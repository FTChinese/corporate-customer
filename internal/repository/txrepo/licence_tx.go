package txrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
)

func (tx TxRepo) RetrieveBaseLicence(r admin.AccessRight) (licence.BaseLicence, error) {
	var bl licence.BaseLicence
	err := tx.Get(
		&bl,
		licence.StmtLockLicence,
		r.RowID,
		r.TeamID)

	if err != nil {
		return licence.BaseLicence{}, err
	}

	return bl, nil
}

// UpdateLicenceStatus after its status changed.
// * when invitation is created;
// * when it is revoked
// * when it is granted
func (tx TxRepo) UpdateLicenceStatus(lic licence.BaseLicence) error {
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
