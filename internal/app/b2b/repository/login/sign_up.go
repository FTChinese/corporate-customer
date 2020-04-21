package login

import (
	"github.com/FTChinese/b2b/internal/app/b2b/model"
	"github.com/FTChinese/b2b/internal/app/b2b/stmt"
)

// SignUp creates a new admin account.
func (env Env) SignUp(s model.AccountInput) error {
	_, err := env.db.NamedExec(stmt.SignUp, s)
	if err != nil {
		logger.WithField("trace", "SignUp").Error(err)
	}

	return nil
}
