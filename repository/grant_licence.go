package repository

import (
	"database/sql"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/jmoiron/sqlx"
)

// FindInvitationByToken tries to find an Invitation by token.
func (env Env) FindInvitationByToken(token string) (admin.Invitation, error) {
	var inv admin.Invitation
	err := env.db.Get(&inv, stmt.InvitationByToken, token)
	if err != nil {
		return inv, err
	}

	return inv, nil
}

// FindInvitedLicence tries to find a licence belong to
// an invitation.
func (env Env) FindInvitedLicence(claims admin.InviteeClaims) (admin.Licence, error) {
	var ls admin.LicenceSchema
	err := env.db.Get(&ls, stmt.InvitedLicence, claims.LicenceID, claims.InvitationID)
	if err != nil {
		return admin.Licence{}, err
	}

	return ls.Licence()
}

// LockLicence locks a licence for update.
func (tx GrantTx) LockLicence(id string) (admin.BaseLicence, error) {
	var l admin.BaseLicence
	err := tx.Get(&l, stmt.LockLicence, id)
	if err != nil {
		return l, err
	}

	return l, nil
}

func (tx GrantTx) InsertMembership(m reader.Membership) error {
	_, err := tx.NamedExec(stmt.InsertMembership, m)

	if err != nil {
		return err
	}

	return nil
}

func (tx GrantTx) UpdateMembership(m reader.Membership) error {
	_, err := tx.NamedExec(stmt.UpdateMembership, m)

	if err != nil {
		return err
	}

	return nil
}

const stmtAcceptInvitation = `
UPDATE b2b.invitation
SET accepted = 1,
	updated_utc = UTC_TIMESTAMP()
WHERE id = ?
LIMIT 1`

// InvitationAccepted turns the accepted field
// of invitation to true.
func (tx GrantTx) InvitationAccepted(id string) error {
	_, err := tx.NamedExec(stmtAcceptInvitation, id)

	if err != nil {
		return err
	}

	return nil
}

const stmtLicenceGranted = `
UPDATE b2b.licence
SET assignee_id = :assignee_id,
	is_active = 1
WHERE id = :licence_id
LIMIT 1`

// LicenceGranted set the assignee_id field
// to user's uuid and turns is_active to true.
func (tx GrantTx) LicenceGranted(l admin.BaseLicence) error {
	_, err := tx.NamedExec(stmtLicenceGranted, l)

	if err != nil {
		return err
	}

	return nil
}
