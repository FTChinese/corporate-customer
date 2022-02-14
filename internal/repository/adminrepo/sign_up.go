package adminrepo

import (
	"github.com/FTChinese/ftacademy/internal/pkg/admin"
)

// SignUp creates a new admin account.
func (env Env) SignUp(a admin.Account) error {
	_, err := env.DBs.Write.NamedExec(admin.StmtCreateAdmin, a)
	if err != nil {
		return err
	}

	return nil
}
