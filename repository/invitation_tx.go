package repository

import (
	"database/sql"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/jmoiron/sqlx"
)

type InvitationTx struct {
	*sqlx.Tx
}

// RetrieveLicence loads the licence when creating an
// invitation for it.
// This retrieves Licence from the point of admin.
func (tx InvitationTx) RetrieveLicence(licenceID, teamID string) (admin.Licence, error) {
	var ls admin.LicenceSchema

	err := tx.Get(&ls, stmt.LockLicenceByID, licenceID, teamID)
	if err != nil {
		return admin.Licence{}, err
	}

	return ls.Licence()
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
