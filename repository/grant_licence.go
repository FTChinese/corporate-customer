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

func (env Env) FindReader(email string) (reader.Reader, error) {
	var r reader.Reader
	err := env.db.Get(&r, stmt.SelectReader, email)
	if err != nil {
		return r, err
	}

	r.Normalize()

	return r, nil
}

func (env Env) CreateReader(s reader.SignUp) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(stmt.CreateReader, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmt.CreateReaderProfile, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(stmt.SaveReaderVrf, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
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
