package login

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
)

// VerifyingAccount retrieves user account by verification token.
func (env Env) VerifyingAccount(token string) (model.Account, error) {
	var a model.Account
	err := env.db.Get(&a, stmt.AccountByVerifier, token)
	if err != nil {
		logger.WithField("trace", "VerifyingAccount").Error(err)
		return a, err
	}

	return a, nil
}

// SetEmailVerified set the verified field to true.
// You should check whether it is already true before
// performing this operation.
func (env Env) SetAccountVerified(account model.Account) error {
	_, err := env.db.Exec(stmt.EmailVerified, account)

	if err != nil {
		logger.WithField("trace", "SetAccountVerified").Error(err)
		return err
	}

	return nil
}
