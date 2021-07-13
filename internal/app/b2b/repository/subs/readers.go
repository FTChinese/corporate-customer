package subs

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
	model2 "github.com/FTChinese/ftacademy/internal/pkg/model"
)

// FindReader by email.
// Deprecated. Use API.
func (env Env) FindReader(email string) (model2.Reader, error) {
	var r model2.Reader
	err := env.dbs.Read.Get(&r, stmt.SelectReader, email)
	if err != nil {
		return r, err
	}

	r.Normalize()

	return r, nil
}

// CreateReader creates new FTC user.
// Deprecated. Use API.
func (env Env) CreateReader(s model2.SignUp) error {
	tx, err := env.dbs.Write.Beginx()
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
