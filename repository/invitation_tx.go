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
	err := tx.Get(&inv, stmt.LockInvitationByID, invitationID, teamID)
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

// UnlinkInvitedLicence removes invitation related
// data from a licence.
func (tx InvitationTx) UnlinkInvitedLicence(licence admin.Licence) error {
	_, err := tx.NamedExec(stmt.RevokeLicenceInvitation, licence)
	if err != nil {
		return err
	}

	return nil
}

func (tx InvitationTx) RevokeLicence(lic admin.BaseLicence) error {
	_, err := tx.NamedExec(stmt.SetLicenceRevoked, lic)
	if err != nil {
		return err
	}

	return nil
}

// The above four methods together revokes an invitation.

// The following methods, plus FindInvitedLicence, performs
// granting licence to a reader.

// FindInvitationByToken retrieves and locks an invitation
// when reader is trying to accept it.
func (tx InvitationTx) FindInvitationByToken(token string) (admin.Invitation, error) {
	var inv admin.Invitation
	err := tx.Get(&inv, stmt.LockInvitationByToken, token)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

// InvitationAccepted marks an invitation as accepted
// so that it cannot be used again.
func (tx InvitationTx) InvitationAccepted(inv admin.Invitation) error {
	_, err := tx.NamedExec(stmt.AcceptInvitation, inv)

	if err != nil {
		return err
	}

	return nil
}

// LicenceGranted links a licence to an assignee.
func (tx InvitationTx) LicenceGranted(l admin.BaseLicence) error {
	_, err := tx.NamedExec(stmt.SetLicenceGranted, l)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveMembership locks a reader's membership row if it present.
func (tx InvitationTx) RetrieveMembership(id string) (reader.Membership, error) {
	var m reader.Membership

	err := tx.Get(&m, stmt.LockMembership, id)
	if err != nil && err != sql.ErrNoRows {
		return m, err
	}

	m.Normalize()

	return m, nil
}

func (tx InvitationTx) InsertMembership(m reader.Membership) error {
	_, err := tx.NamedExec(stmt.InsertMembership, m)

	if err != nil {
		return err
	}

	return nil
}

func (tx InvitationTx) UpdateMembership(m reader.Membership) error {
	_, err := tx.NamedExec(stmt.UpdateMembership, m)

	if err != nil {
		return err
	}

	return nil
}
