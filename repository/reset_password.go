package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// InsertPasswordResetToken generates a new token
// to help resetting password.
const stmtInsertPwResetToken = `
INSERT INTO b2b.password_reset
SET token = UNHEX(:token),
	email = :email,
	created_utc = UTC_TIMESTAMP()`

// SavePasswordResetter
func (env Env) SavePasswordResetter(pr admin.PasswordResetter) error {
	_, err := env.db.NamedExec(stmtInsertPwResetToken, pr)

	if err != nil {
		return err
	}

	return nil
}

// RetrievePasswordReset retrieves an account for a
// password reset token.
const stmtPwResetAccount = stmt.AccountBase + `
FROM b2b.password_reset AS r
	INNER JOIN b2b.admin AS a
	ON r.email = a.email
WHERE r.token = UNHEX(?)
	AND r.is_used = 0
	 DATE_ADD(r.created_utc, INTERVAL r.expires_in SECOND) > UTC_TIMESTAMP() 
ORDER BY r.created_utc DESC
LIMIT 1`

func (env Env) PasswordResettingAccount(token string) (admin.Account, error) {
	var a admin.Account
	if err := env.db.Get(&a, stmtPwResetAccount, token); err != nil {
		return a, err
	}

	return a, nil
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
