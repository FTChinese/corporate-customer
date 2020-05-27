package login

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
)

// FindPwResetAccount tries to find the account by email
// user submitted to request a password resetting letter.
// This email is used to receive a password reset letter.
func (env Env) FindPwResetAccount(email string) (model.Account, error) {
	var a model.Account
	err := env.db.Get(&a, stmt.AccountByEmail, email)

	if err != nil {
		logger.WithField("trace", "FindPwResetAccount").Error(err)
		return a, err
	}

	return a, nil
}

// SavePwResetToken saves password reset token together the the
// email linked to it.
func (env Env) SavePwResetToken(pr model.AccountInput) error {
	_, err := env.db.NamedExec(stmt.InsertPwResetToken, pr)

	if err != nil {
		return err
	}

	return nil
}

// FIndAccountByPwToken retrieves an account for a
// password reset token.
func (env Env) FindAccountByPwToken(token string) (model.Account, error) {
	var a model.Account
	if err := env.db.Get(&a, stmt.AccountForPwReset, token); err != nil {
		logger.WithField("trace", "FindAccountByPwToken").Error(err)
		return a, err
	}

	return a, nil
}

// ResetPassword resets admin's password using its id.
func (env Env) ResetPassword(credentials model.AccountInput) error {
	_, err := env.db.NamedExec(stmt.UpdatePassword, credentials)
	if err != nil {
		logger.WithField("trace", "ResetPassword")
		return err
	}

	return nil
}

// RemovePwResetToken disables a token so that is cannot be
// used again.
func (env Env) RemovePwResetToken(token string) error {
	_, err := env.db.Exec(stmt.DeactivatePwResetToken, token)
	if err != nil {
		return err
	}

	return nil
}
