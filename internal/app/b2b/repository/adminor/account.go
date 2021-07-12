package adminor

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
)

func (env Env) BaseAccountByID(id string) (admin.BaseAccount, error) {
	var a admin.BaseAccount
	err := env.DBs.Read.Get(&a, admin.StmtBaseAccountByID, id)
	if err != nil {
		return admin.BaseAccount{}, err
	}

	return a, nil
}

func (env Env) BaseAccountByEmail(email string) (admin.BaseAccount, error) {
	var a admin.BaseAccount
	err := env.DBs.Read.Get(&a, admin.StmtBaseAccountByEmail, email)
	if err != nil {
		return admin.BaseAccount{}, err
	}

	return a, nil
}

// UpdateName allows admin to change display name.
func (env Env) UpdateName(account admin.BaseAccount) error {
	_, err := env.DBs.Write.NamedExec(admin.StmtUpdateName, account)

	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates reader's password.
// This is used both by resetting password if forgotten and updating password after logged in.
func (env Env) UpdatePassword(p input.PasswordUpdateParams) error {

	_, err := env.DBs.Write.NamedExec(admin.StmtUpdatePassword,
		p)

	if err != nil {
		return err
	}

	return nil
}
