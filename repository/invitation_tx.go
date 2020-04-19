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

// SetLicenceInvited links a licence to an invitation.
func (tx InvitationTx) SetLicenceInvited(lic admin.BaseLicence) error {
	_, err := tx.NamedExec(stmt.SetLicenceInvited, lic)
	if err != nil {
		return err
	}

	return nil
}

// The above methods creates an invitation together.

// RetrieveInvitation when admin wants to revoke it.
// This retrieves licence from the point of Invitation.
func (tx InvitationTx) RetrieveInvitation(invitationID, teamID string) (admin.Invitation, error) {
	var inv admin.Invitation
	err := tx.Get(&inv, stmt.LockInvitation, invitationID, teamID)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

// FindInvitedLicence retrieves the licence belonging
// to an invitation.
// A licence could have multiple invitations. Since a licence
// is only linked to the last invitation created under it,
// with licence id only you could not tell whether this
// licence is still linked to this invitation. Therefore
// we use the last_invitation_id column to load it.
func (tx InvitationTx) FindInvitedLicence(inv admin.Invitation) (admin.Licence, error) {
	var ls admin.LicenceSchema
	err := tx.Get(&ls, stmt.LockInvitedLicence, inv.LicenceID, inv.ID)
	if err != nil {
		return admin.Licence{}, err
	}

	return ls.Licence()
}

// RevokeInvitation marks an invitation as revoked.
// The corresponding licence should also remove any traces
// linking to this invitation.
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
