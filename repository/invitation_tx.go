package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/jmoiron/sqlx"
)

type InvitationTx struct {
	*sqlx.Tx
}

// RetrieveLicence finds the licence used in an invitation
func (tx InvitationTx) RetrieveLicence(licenceID, teamID string) (admin.Licence, error) {
	var ls admin.LicenceSchema

	err := tx.Get(&ls, stmt.LockLicence, licenceID, teamID)
	if err != nil {
		return admin.Licence{}, err
	}

	return ls.Licence()
}

// SetLicenceInvited changes a licence status and set the invitation column.
func (tx InvitationTx) SetLicenceInvited(lic admin.BaseLicence) error {
	_, err := tx.NamedExec(stmt.SetLicenceInvited, lic)
	if err != nil {
		return err
	}

	return nil
}

// SaveInvitation insert a new row into invitation table.
func (tx InvitationTx) SaveInvitation(inv admin.Invitation) error {
	_, err := tx.NamedExec(stmt.CreateInvitation, inv)

	if err != nil {
		return err
	}

	return nil
}

func (tx InvitationTx) RetrieveInvitation(invitationID, teamID string) (admin.Invitation, error) {
	var inv admin.Invitation
	err := tx.Get(&inv, stmt.LockInvitation, invitationID, teamID)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

func (tx InvitationTx) RevokeInvitation(inv admin.Invitation) error {
	_, err := tx.NamedExec(stmt.RevokeInvitation, inv)
	if err != nil {
		return err
	}

	return nil
}

func (tx InvitationTx) UnlinkLicenceInvitation(licence admin.Licence) error {
	_, err := tx.NamedExec(stmt.RevokeLicenceInvitation, licence)
	if err != nil {
		return err
	}

	return nil
}
