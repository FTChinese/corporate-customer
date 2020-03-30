package repository

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

const stmtAccountByID = stmt.AccountBase + `
FROM b2b.admin AS a
WHERE a.id = ?
LIMIT 1`

// AccountByID retrieves user account by id
func (env Env) AccountByID(id string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmtAccountByID, id)
	if err != nil {
		return a, err
	}

	return a, nil
}

const stmtAccountByVerifier = stmt.AccountBase + `
FROM b2b.admin AS a
WHERE a.vrf_token = ?
LIMIT 1`

// AccountByVerifier retrieves user account by verification token.
func (env Env) AccountByVerifier(token string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmtAccountByVerifier, token)
	if err != nil {
		return a, err
	}

	return a, nil
}

// SetEmailVerified set the verified column to true
// for a verification token.
const stmtEmailVerified = `
UPDATE b2b.admin
	SET verified = 1
WHERE vrf_token = ?
LIMIT 1`

// SetEmailVerified set the verified field to true.
// You should check whether it is already true before
// performing this operation.
func (env Env) SetEmailVerified(token string) error {
	_, err := env.db.Exec(stmtEmailVerified, token)

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

const stmtUpdatePassword = `
UPDATE b2b.admin
SET password_sha2 = UNHEX(SHA2(:password, 256))
	updated_utc = UTC_TIMESTAMP()
WHERE id = :amin_id
LIMIT 1`

func (env Env) UpdatePassword(c admin.Credentials) error {
	_, err := env.db.NamedExec(stmtUpdatePassword, c)
	if err != nil {
		return err
	}

	return nil
}

// ReGeneratedVrfToken regenerate verification token upon request.
// An email should also be sent.
const stmtReGenerateVrfToken = `
UPDATE b2b.admin
SET token = UNHEX(:oken),
	updated_utc = UTC_TIMESTAMP()
WHERE id = :admin_id
LIMIT1`

// RegenerateVerifier re-generate a verification token
// upon user request. If a account is already verified,
// the request should be ignored.
func (env Env) RegenerateVerifier(v admin.Verifier) error {
	_, err := env.db.NamedExec(stmtReGenerateVrfToken, v)

	if err != nil {
		return err
	}

	return nil
}
