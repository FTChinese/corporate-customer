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

// TakeSnapshot backs up a membership before
// modifying it.
func (env Env) TakeSnapshot(snp reader.MemberSnapshot) error {
	_, err := env.db.NamedExec(stmt.TakeSnapshot, snp)

	if err != nil {
		return err
	}

	return nil
}
