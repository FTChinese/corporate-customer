package login

import (
	"github.com/FTChinese/ftacademy/internal/app/b2b/model"
	"github.com/FTChinese/ftacademy/internal/app/b2b/stmt"
)

// Login verifies user email + password.
func (env Env) Login(in model.AccountInput) (model.Passport, error) {
	var pp model.Passport

	err := env.db.Get(&pp, stmt.Login, in.Email, in.Password)

	if err != nil {
		logger.WithField("trace", "Login").Error(err)
		return pp, err
	}

	return pp, nil
}
