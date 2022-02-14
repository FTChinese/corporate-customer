package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
)

// Authenticate verifies email+password are correct.
// If the returned error is sql.ErrorNoRow, it indicates the
// email does not exists.
// If no error returned, the AuthResult.PasswordMatched
// field indicates whether the password is correct.
func (env Env) Authenticate(params input.Credentials) (admin.AuthResult, error) {
	var r admin.AuthResult
	err := env.DBs.Read.Get(&r,
		admin.StmtVerifyPasswordByEmail,
		params.Password,
		params.Email)

	if err != nil {
		return r, err
	}

	return r, nil
}

// VerifyPassword when user is trying to change it.
func (env Env) VerifyPassword(params input.PasswordUpdateParams) (admin.AuthResult, error) {
	var r admin.AuthResult
	err := env.DBs.Read.Get(
		&r,
		admin.StmtVerifyPasswordByID,
		params.Old,
		params.ID)

	if err != nil {
		return admin.AuthResult{}, err
	}

	return r, nil
}
