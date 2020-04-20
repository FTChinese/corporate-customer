package setting

import (
	"github.com/FTChinese/b2b/models/admin"
	"github.com/FTChinese/b2b/repository/stmt"
)

// PassportByAdminID retrieves admin's account and team data
// by admin id.
func (env Env) LoadPassport(adminID string) (admin.Passport, error) {
	var a admin.Passport
	if err := env.db.Get(&a, stmt.PassportByAdminID, adminID); err != nil {
		logger.WithField("trace", "LoadPassport").Error()
		return a, err
	}

	return a, nil
}

// Account retrieves user account by id
func (env Env) Account(adminID string) (admin.Account, error) {
	var a admin.Account
	err := env.db.Get(&a, stmt.AccountByID, adminID)
	if err != nil {
		logger.WithField("trace", "RetrieveAccount").Error(err)
		return a, err
	}

	return a, nil
}

// Profile loads admin's full account data.
func (env Env) Profile(id string) (admin.Profile, error) {
	var p admin.Profile
	err := env.db.Get(&p, stmt.AdminProfile, id)
	if err != nil {
		return p, err
	}

	return p, nil
}

// RegenerateVerifier re-generate a verification token
// upon user request. If a account is already verified,
// the request should be ignored.
func (env Env) RegenerateVerifier(verifier admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.ReGenerateVrfToken, verifier)

	if err != nil {
		return err
	}

	return nil
}

// UpdateName allows admin to change display name.
func (env Env) UpdateName(account admin.Account) error {
	_, err := env.db.NamedExec(stmt.UpdateName, account)

	if err != nil {
		return err
	}

	return nil
}

// PasswordMatched checks whether user's current password is correct.
func (env Env) PasswordMatched(input admin.AccountInput) (bool, error) {
	var ok bool
	if err := env.db.Get(&ok, stmt.PasswordMatched, input.OldPassword, input.ID); err != nil {
		return ok, err
	}

	return ok, nil
}

// UpdatePassword allows user to change password.
func (env Env) UpdatePassword(input admin.AccountInput) error {
	_, err := env.db.NamedExec(stmt.UpdatePassword, input)
	if err != nil {
		return err
	}

	return nil
}
