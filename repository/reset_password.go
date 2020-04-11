package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

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

// InsertPasswordResetToken generates a new token
// to help resetting password.
const stmtInsertPwResetToken = `
INSERT INTO b2b.password_reset
SET email = :email,
	token = UNHEX(:token),
	is_used = 0,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()
ON DUPLICATE KEY UPDATE
	token = UNHEX(:token),
	is_used = 0,
	updated_utc = UTC_TIMESTAMP()`

// SavePasswordResetter
func (env Env) SavePasswordResetter(pr admin.AccountInput) error {
	_, err := env.db.NamedExec(stmtInsertPwResetToken, pr)

	if err != nil {
		return err
	}

	return nil
}

// DeactivatePasswordResetToken set a token's is_used
// column to true so that it won't be retrieved in the future.
const stmtDeactivatePwResetToken = `
UPDATE b2b.password_reset
	SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`

func (env Env) RemovePasswordResetToken(token string) error {
	_, err := env.db.Exec(stmtDeactivatePwResetToken, token)
	if err != nil {
		return err
	}

	return nil
}
