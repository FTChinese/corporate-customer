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
func (tx InvitationTx) RetrieveLicence(input admin.InvitationInput) (admin.LicenceSchema, error) {
	var lic admin.LicenceSchema

	err := tx.Get(&lic, stmt.LockLicence, input.LicenceID, input.TeamID)
	if err != nil {
		_ = tx.Rollback()
		return lic, err
	}

	return lic, nil
}

// UpdateLicence changes a licence status and set the invitation
// column.
func (tx InvitationTx) UpdateLicence(lic admin.LicenceSchema) error {
	_, err := tx.NamedExec(stmt.SetLicenceInvited, lic)
	if err != nil {
		_ = tx.Rollback()
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
