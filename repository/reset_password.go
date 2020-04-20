package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// AccountByEmail loads admin account by email.
// This email is used to receive a password reset letter.
func (env Env) AccountByEmail(email string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmt.AccountByEmail, email)

	if err != nil {
		return a, err
	}

	return a, nil
}

// PasswordResettingAccount retrieves an account for a
// password reset token.
func (env Env) AccountToResetPassword(token string) (admin.Account, error) {
	var a admin.Account
	if err := env.db.Get(&a, stmt.AccountForPwReset, token); err != nil {
		return a, err
	}

	return a, nil
}

// SavePasswordResetter
func (env Env) SavePasswordResetter(pr admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.InsertPwResetToken, pr)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) RemovePasswordResetToken(token string) error {
	_, err := env.db.Exec(stmt.DeactivatePwResetToken, token)
	if err != nil {
		return err
	}

	return nil
}
