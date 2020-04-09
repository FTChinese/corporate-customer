package repository

import (
	"database/sql"
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/models/reader"
	"github.com/FTChinese/b2b/repository/stmt"
	"github.com/jmoiron/sqlx"
)

// GrantTx splits a series of db transaction into smaller functions.
// When grating a licence, we need to lock reader's current membership,
// the invitation, and licence.
// We also need to backup current membership if exists.
type GrantTx struct {
	*sqlx.Tx
}

const stmtMembership = `
SELECT ` + stmt.MembershipSelectCols + `
FROM premium.ftc_vip AS m
WHERE vip_id = ?
LIMIT 1
FOR UPDATE`

// RetrieveMembership locks a reader's membership row if it present.
func (tx GrantTx) LockMembership(id string) (reader.Membership, error) {
	var m reader.Membership

	err := tx.Get(&m, stmtMembership)
	if err != nil && err != sql.ErrNoRows {
		return m, err
	}

	m.Normalize()

	return m, nil
}

// LockInvitation locks an invitation for update.
func (tx GrantTx) LockInvitation(id string) (admin.Invitation, error) {
	var i admin.Invitation
	err := tx.Get(&i, stmt.LockInvitation, id)
	if err != nil {
		return i, err
	}

	return i, nil
}

// LockLicence locks a licence for update.
func (tx GrantTx) LockLicence(id string) (admin.Licence, error) {
	var l admin.Licence
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
func (tx GrantTx) LicenceGranted(l admin.Licence) error {
	_, err := tx.NamedExec(stmtLicenceGranted, l)

	if err != nil {
		return err
	}

	return nil
}
