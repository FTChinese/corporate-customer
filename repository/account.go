package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// AccountByID retrieves user account by id
func (env Env) AccountByID(id string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmt.AccountByID, id)
	if err != nil {
		return a, err
	}

	return a, nil
}

// AccountByVerifier retrieves user account by verification token.
func (env Env) AccountByVerifier(token string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmt.AccountByVerifier, token)
	if err != nil {
		return a, err
	}

	return a, nil
}

// AdminProfile loads admin's full account data.
func (env Env) AdminProfile(id string) (admin.Profile, error) {
	var p admin.Profile
	err := env.db.Get(&p, stmt.AdminProfile, id)
	if err != nil {
		return p, err
	}

	return p, nil
}

// SetEmailVerified set the verified column to true
// for a verification token.
const stmtEmailVerified = `
UPDATE b2b.admin
	SET verified = 1
WHERE id = :admin_id
LIMIT 1`

// SetEmailVerified set the verified field to true.
// You should check whether it is already true before
// performing this operation.
func (env Env) SetEmailVerified(account admin.Account) error {
	_, err := env.db.Exec(stmtEmailVerified, account)

	if err != nil {
		return err
	}

	return nil
}

// ReGeneratedVrfToken regenerate verification token upon request.
// An email should also be sent.
const stmtReGenerateVrfToken = `
UPDATE b2b.admin
SET token = UNHEX(:token),
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT1`

// RegenerateVerifier re-generate a verification token
// upon user request. If a account is already verified,
// the request should be ignored.
func (env Env) RegenerateVerifier(verifier admin.AccountInput) error {
	_, err := env.db.NamedExec(stmtReGenerateVrfToken, verifier)

	if err != nil {
		return err
	}

	return nil
}

const stmtUpdateName = `
UPDATE b2b.admin
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT 1`

func (env Env) UpdateName(account admin.Account) error {
	_, err := env.db.NamedExec(stmtUpdateName, account)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) PasswordMatched(input admin.AccountInput) (bool, error) {
	var ok bool
	if err := env.db.Get(&ok, stmt.PasswordMatched, input.OldPassword, input.ID); err != nil {
		return ok, err
	}

	return ok, nil
}

const stmtUpdatePassword = `
UPDATE b2b.admin
SET password_sha2 = UNHEX(SHA2(:password, 256))
	updated_utc = UTC_TIMESTAMP()
WHERE id = :amin_id
LIMIT 1`

func (env Env) UpdatePassword(input admin.AccountInput) error {
	_, err := env.db.NamedExec(stmtUpdatePassword, input)
	if err != nil {
		return err
	}

	return nil
}
